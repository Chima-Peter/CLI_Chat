# Feature IN-02 — CLI Server

**Slug:** cli-server  
**Output file:** IN-02-cli-server.md

## Summary

The CLI Server (`IN-02`) serves as the entry point for the CLI Chat server application. It is responsible for initializing the server instance, configuring TLS security, and binding the TCP listener to a designated port.

## Scope

- **In scope**: Server startup, TLS configuration, TCP listener setup, and incoming connection acceptance.
- **Out of scope**: Business logic for chat rooms, message handling, and user management (delegated to the `server` package).

## Functional behaviour

When executed, the CLI Server:
1. Displays a bootstrap banner with server information and supported client commands (`cmd/server/main.go:13-28`).
2. Initializes the server components (`cmd/server/main.go:33`).
3. Sets up a TLS configuration for secure communication (`cmd/server/main.go:35-38`).
4. Listens for incoming TCP connections on a configurable port (`cmd/server/main.go:40-45`), defaulting to `8080` if the `PORT` environment variable is not set.
5. Accepts incoming connections in a loop and spawns a new goroutine for each connection to handle them concurrently (`cmd/server/main.go:55-62`).

## Technical design

The server is implemented in Go as a console application. It relies on:
- **`crypto/tls`**: For establishing secure TCP connections.
- **`server` package**: For the core server logic, including `InitServer()` and `HandleConn()`.
- **`tls_config` package**: For generating development-ready TLS configurations.

Configuration is driven by environment variables, specifically `PORT`.

## Entry points

- `main()` function in `cmd/server/main.go` serves as the primary entry point for the CLI server binary.

## Dependencies

- **Standard Library**: `crypto/tls`, `fmt`, `log`, `os`.
- **Internal Packages**:
  - `github.com/chima/CLI_Chat/server`
  - `github.com/chima/CLI_Chat/tls_config`

## Open questions

None.

## Revision

- Initial draft: scope and behaviour from implementation reads.

## Evidence index

- CLI Server entry point implementation: `cmd/server/main.go:1-64`
- Bootstrap banner: `cmd/server/main.go:13-28`
- Server initialization and connection handling: `cmd/server/main.go:33-62`
- Port configuration: `cmd/server/main.go:40-43`

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [CLI Server Main Entry](./IN-02_IC-01-cli-server-main-entry.md)
- [Server State Initialization](./IN-02_IC-02-server-state-initialization.md)
- [Connection Handling](./IN-02_IC-03-connection-handling.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



