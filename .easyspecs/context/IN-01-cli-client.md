# Feature IN-01 — CLI Client

**Slug:** cli-client  
**Output file:** IN-01-cli-client.md

## Summary

The CLI Client (IN-01) is the primary command-line interface entry point for the application. It provides the user with an interactive terminal session to interact with the server.

## Scope

- **In-scope:**
    - CLI entry point initialization.
    - Establishing a TLS-encrypted connection to the server.
    - Providing an interactive command loop via `readline`.
    - Session management.
- **Out-of-scope:**
    - Server-side implementation (handled by IN-02).
    - TLS certificate generation (handled by IN-03).
    - Authentication logic itself (delegated to the server once connected).

## Functional behaviour

The CLI client functions as a terminal-based application. Upon execution, it:
1. Performs a bootstrap process.
2. Reads the `SERVER_URL` environment variable to determine the target server, defaulting to `localhost:8080`.
3. Establishes a TLS-encrypted TCP connection to the specified server with insecure skip verify enabled.
4. Initializes an interactive `readline` session.
5. Invokes the session loop where the user can interact with the system.
6. Handles graceful disconnection.

## Technical design

The implementation relies on:
- `cmd/client/main.go`: The main entry point that calls `client.Connect()`.
- `client/client.go`: Handles the connection establishment, `readline` configuration, and the session loop.
- `crypto/tls`: For encrypted communication.
- `github.com/chzyer/readline`: For providing a feature-rich CLI interface.

## Entry points

- **CLI:** `cmd/client/main.go`

## Dependencies

- **Go Standard Library:** `crypto/tls`, `fmt`, `os`
- **External:** `github.com/chzyer/readline` (for input handling)
- **Infrastructure:** `IN-02` (CLI Server) for target connections.

## Open questions

None.

## Revision

- Initial draft: scope and behaviour from implementation reads.

## Evidence index

- `cmd/client/main.go:1-9`: CLI entry point definition.
- `client/client.go:11-41`: Implementation of `Connect()` function, TLS setup, and session initialization.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [CLI client main entry](./IN-01_IC-01-ic01.md)
- [Client session management](./IN-01_IC-02-ic02.md)
- [Client message and command handling](./IN-01_IC-03-ic03.md)
- [Client file transfer](./IN-01_IC-04-ic04.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



