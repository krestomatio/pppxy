# Repository guidance

Follow this repository's documentation and any more specific instructions for
the files being changed.

## Repository scope

- PPPxy is the Go TCP proxy used in the private-cloud load-balancer path to
  prepend PROXY protocol v1/v2 headers before forwarding connections.
- `cmd/pppxy/` is the executable entry point; `pkg/pppxy/` owns configuration and
  connection behavior; `config.yaml` documents the runtime schema.
- Preserve byte-level PROXY protocol correctness, connection/timeout behavior,
  multi-listener isolation, graceful shutdown, non-root operation, and backward
  compatibility of config keys.

## Validation

- Run `go test ./...`, `go vet ./...`, and the configured Go linter/`make go-lint`.
  Add focused tests for parsing and wire behavior when changing the proxy path.
- For integration checks, use loopback listeners and synthetic backends; do not
  bind production ports or forward to real infrastructure.
