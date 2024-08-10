# PPPxy

**PPPxy (PROXY Protocol Proxy)** is a lightweight proxy application designed to transparently add the [PROXY protocol](https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) header to TCP connections. The app supports handling multiple proxy servers, each with its own configuration.

## Features

- Supports Proxy Protocol versions 1 and 2.
- Configurable via a YAML file.

## Installation

### Binary

To install PPPxy, clone the repository and build the binary:

```bash
git clone <repository-url>
cd pppxy
go build -o pppxy cmd/pppxy/main.go
```

### Container
To run as a container:
```bash
# create config.yaml
# CAP_NET_BIND_SERVICE is for binding to privileged ports
# or docker
podman run --rm -d \
  --cap-add=CAP_NET_BIND_SERVICE \
  --net=host \
  --user 1001 \
  --volume "${PWD}/config.yaml:/etc/pppxy/config.yaml:ro" \
  --name pppxy quay.io/krestomatio/pppxy
```

### Systemd
To install and run as a systemd service:
```bash
mkdir -p /etc/pppxy
# modify the configuration file as needed
cat << EOF > /etc/pppxy/config.yaml
pppxy_group:
  - listen_addr: ":11443"
    backend_addr: "192.168.1.2:22443"
    proxy_protocol_version: 1
  - listen_addr: ":33443"
    backend_addr: "192.168.1.3:44443"
    proxy_protocol_version: 2
EOF
chown 1001:1001 /etc/pppxy/config.yaml
chcon -t container_file_t /etc/pppxy/config.yaml
podman create --rm --restart on-failure \
  --stop-timeout 30 \
  --cap-add=CAP_NET_BIND_SERVICE \
  --net=host \
  --user 1001 \
  --volume /etc/pppxy/config.yaml:/etc/pppxy/config.yaml:ro \
  --name pppxy quay.io/krestomatio/pppxy:0.0.1
podman generate systemd --new \
  --restart-sec 15 \
  --start-timeout 180 \
  --stop-timeout 30 \
  --name pppxy > /etc/systemd/system/pppxy.service
systemctl daemon-reload
systemctl enable --now pppxy.service
```

## Usage

Run the `pppxy` application with the desired configuration file:

```bash
./pppxy --help
Usage of ./pppxy:
  -config string
        Path to configuration file (default "/etc/pppxy/config.yaml")
  -log-level string
        Log level (debug, info, warn, error) (default "info")
```

## Configuration

The application is configured via a YAML file. Below is an example configuration:

```yaml
pppxy_group:
  - listen_addr: ":11443"
    backend_addr: "192.168.1.2:22443"
    proxy_protocol_version: 1
  - listen_addr: ":33443"
    backend_addr: "192.168.1.3:44443"
    proxy_protocol_version: 2
```

### Configuration Options

- **`listen_addr`**: The address and port where the proxy listens for incoming connections.
- **`backend_addr`**: The backend server address and port where the connection will be forwarded.
- **`proxy_protocol_version`**: The Proxy Protocol version to use (1 or 2).

## Log Levels

- **debug**: Detailed logs for debugging.
- **info**: General operational logs.
- **warn**: Logs that indicate potential issues.
- **error**: Logs that indicate failures.

## Krestomatio

This project is part of open source contribution at Krestomatio, a service offering [managed Moodle™ e-learning platforms](https://krestomatio.com).
