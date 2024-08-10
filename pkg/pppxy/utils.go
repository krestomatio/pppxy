package pppxy

import (
	"errors"
	"net"
	"syscall"
)

// isSocketNotConnectedError checks if the error is due to a socket not being connected.
func isSocketNotConnectedError(err error) bool {
	var oerr *net.OpError
	return errors.As(err, &oerr) && errors.Is(err, syscall.ENOTCONN)
}

// isConnectionResetDuringRead checks if the error is a connection reset during a read operation.
func isConnectionResetDuringRead(err error) bool {
	var oerr *net.OpError
	return errors.As(err, &oerr) && oerr.Op == "read" && errors.Is(err, syscall.ECONNRESET)
}

// isUseOfClosedNetworkConnection checks if the error is due to the use of a closed network connection.
func isUseOfClosedNetworkConnection(err error) bool {
	var opErr *net.OpError
	return errors.As(err, &opErr) && opErr.Err.Error() == "use of closed network connection"
}
