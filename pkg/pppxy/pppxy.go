package pppxy

import (
	"fmt"
	"io"
	"net"

	"github.com/go-logr/logr"
	"github.com/pires/go-proxyproto"
)

// PPPxy represents a proxy protocol Proxy
type PPPxy struct {
	Config   PPPxyConfig
	Log      logr.Logger
	dialer   net.Dialer
	listener net.Listener
}

// WriteCloser describes a net.Conn with a CloseWrite method.
type WriteCloser interface {
	net.Conn
	CloseWrite() error
}

// NewPPPxy creates a new proxy protocol proxy instance.
func NewPPPxy(config PPPxyConfig, log logr.Logger) *PPPxy {
	return &PPPxy{
		Config: config,
		Log:    log,
	}
}

// Start initializes the proxy protocol proxy and starts listening for connections.
func (p *PPPxy) Start() error {
	if err := p.Listen(); err != nil {
		return fmt.Errorf("failed to start proxy protocol proxy: %w", err)
	}

	go p.Handle()
	return nil
}

// Listen starts the proxy protocol proxy and binds to the specified address.
func (p *PPPxy) Listen() error {
	listener, err := net.Listen("tcp", p.Config.ListenAddr)
	if err != nil {
		return err
	}
	p.listener = listener

	p.Log.Info("Listening", "address", p.Config.ListenAddr, "forwarding to", p.Config.BackendAddr)
	return nil
}

// Handle accepts and processes incoming client connections.
func (p *PPPxy) Handle() {
	defer p.listener.Close()

	for {
		connClient, err := p.listener.Accept()
		if err != nil {
			// Check if the error is due to the listener being closed and not isUseOfClosedNetworkConnection
			if isUseOfClosedNetworkConnection(err) {
				p.Log.Info("Listener closed, stopping server", "address", p.Config.ListenAddr)
				return
			}
			p.Log.Error(err, "Error accepting connection")
			continue
		}
		p.Log.V(1).Info("Accepted connection", "client", connClient.RemoteAddr().String())
		go p.handleClientConnection(connClient.(WriteCloser))
	}
}

// handleClientConnection handles a client connection by forwarding traffic to the backend server.
func (p *PPPxy) handleClientConnection(connClient WriteCloser) {
	defer connClient.Close()

	connBackend, err := p.dialBackend()
	if err != nil {
		p.Log.Error(err, "Failed to connect to backend server")
		return
	}
	defer connBackend.Close()

	if err := p.sendProxyProtocolHeader(connClient, connBackend); err != nil {
		p.Log.Error(err, "Failed to write proxy protocol header")
		return
	}

	p.forwardTraffic(connClient, connBackend)
}

// sendProxyProtocolHeader sends a proxy protocol header to the backend server.
func (p *PPPxy) sendProxyProtocolHeader(connClient net.Conn, connBackend WriteCloser) error {
	sourceAddr := connClient.RemoteAddr()
	destAddr := connClient.LocalAddr()
	header := proxyproto.HeaderProxyFromAddrs(byte(p.Config.ProxyProtocolVersion), sourceAddr, destAddr)
	if _, err := header.WriteTo(connBackend); err != nil {
		return err
	}

	p.Log.V(1).Info("Sent proxy protocol header", "header", header)
	return nil
}

// dialBackend establishes a connection to the backend server.
func (p *PPPxy) dialBackend() (WriteCloser, error) {
	conn, err := p.dialer.Dial("tcp", p.Config.BackendAddr)
	if err != nil {
		return nil, err
	}
	return conn.(WriteCloser), nil
}

// forwardTraffic forwards traffic between the client and backend connections.
func (p *PPPxy) forwardTraffic(connClient, connBackend WriteCloser) {
	errChan := make(chan error, 2)
	go p.copyConnection(connBackend, connClient, errChan)
	go p.copyConnection(connClient, connBackend, errChan)

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			if isConnectionResetDuringRead(err) {
				p.Log.V(1).Info("Connection reset during read", "error", err)
			} else {
				p.Log.Error(err, "Connection error")
			}
		}
	}
}

// copyConnection copies data between two connections and handles connection closure.
func (p *PPPxy) copyConnection(dst, src WriteCloser, errCh chan error) {
	_, err := io.Copy(dst, src)
	errCh <- err

	if err := dst.CloseWrite(); err != nil && !isSocketNotConnectedError(err) {
		p.Log.Error(err, "Failed to close write on connection")
	}
}

// Close shuts down the proxy protocol proxy.
func (p *PPPxy) Close() error {
	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}
