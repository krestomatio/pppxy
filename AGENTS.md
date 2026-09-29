# Repository guidance

Before working, read these optional instruction files in order,
resolving these paths from this repository's root:

1. `.agents/organization/AGENTS.md`
2. `.agents/workspace/AGENTS.md`

Read each file if it resolves to a readable regular file.
Skip absent files; report broken or unreadable links.
Read each resolved file only once to avoid duplicate loading and cycles.

Resolve references inside imported files relative to their real target
directory after following symlinks.

Apply organization guidance, then workspace guidance, then the repository
instructions below. More specific applicable instructions take precedence.

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
