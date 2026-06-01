# Method ME-05 — Connect

**Service:** SV-02 · **File:** SV-02_ME-05-connect.md

## Summary

The `Connect` method initializes the client connection to the server, establishes a TLS-encrypted TCP connection, sets up the interactive command interface, and initiates the client-server session.

## Operation

`Connect()` function in `client/client.go`. This is the entry point for starting the client application and connecting to the server.

## Request / inputs

- **Environment Variable:** `SERVER_URL` (optional). If not set, defaults to `localhost:8080`.

## Response / outputs

- **Console Output:**
  - Connection status messages (success or failure).
  - Terminal prompt for user interaction via `readline`.

## Auth and permissions

- The connection is established via TLS (`tls.Dial`) with `InsecureSkipVerify: true`, effectively disabling server certificate verification.

## Idempotency and concurrency

- The operation is designed to be the single entry point for establishing a session. It is not designed to be called concurrently from the same process.

## Errors

- **Connection Error:** If `tls.Dial` fails, the error is handled by printing a user-friendly message to the console ("Server currently down...") and returning from the function.
- **Readline Error:** If `readline.NewEx` fails, the client panics.

## Implementation notes

- The implementation relies on standard Go TLS for networking and `chzyer/readline` for CLI interaction.
- The connection and the readline interface are closed on function exit using `defer`.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `client/client.go:11-41`: Implementation of the `Connect` method, including server URL retrieval, TLS connection setup, error handling, and session initiation.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./SV-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



