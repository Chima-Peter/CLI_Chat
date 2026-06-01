# Method ME-19 — handleFilePortListening

**Service:** SV-01 · **File:** SV-01_ME-19-handle-file-port-listening.md

## Summary

Configures the file listening host and port for a connected client, enabling file transfers, and subsequently dispatches any pending file transfer requests for that client.

## Operation

`handleFilePortListening(receiver *client, payload json.RawMessage)`

Processes a JSON payload containing `host` and `port` information, validates the provided values, updates the client's file listening configuration, and attempts to fulfill any queued file transfer requests for the client.

## Request / inputs

The method receives a `json.RawMessage` payload expected to be a JSON object:

- `host` (string): The hostname or IP address where the client is listening for file transfer connections.
- `port` (string): The port number where the client is listening, which must be a valid integer between 1 and 65535 inclusive.

## Response / outputs

No direct return value. The method modifies the internal state of the `receiver` client (sets `fileListenHost` and `fileListenPort`).

## Auth and permissions

The method is intended to be called by a connected client. Authentication is implicitly required for the `receiver` client to be established in the server's session management.

## Idempotency and concurrency

- **Concurrency:** The updates to `receiver.fileListenHost` and `receiver.fileListenPort` are protected by `receiver.mu.Lock()` and `receiver.mu.Unlock()` to ensure thread-safe state modification.
- **Idempotency:** Subsequent calls with the same host and port are generally safe, as they simply overwrite the existing configuration and re-trigger pending request dispatching.

## Errors

- Invalid `port` (not an integer, or out of range [1, 65535]): The client receives an error message indicating an invalid file port.
- Missing `host`: The client receives an error message requiring the file port host.
- Errors are communicated back to the client using `receiver.err(error)`.

## Implementation notes

- The `port` is parsed from the `meta.Port` string using `strconv.Atoi` after trimming whitespace.
- The `host` is trimmed using `strings.TrimSpace`.
- The method triggers `s.dispatchPendingFileRequests(receiver)` after successfully configuring the listener, ensuring immediate processing of any waiting file transfers.

## Revision

- Initial draft: implementation contract and handler logic documented.

## Evidence index

- `server/server.go:443-468`: Implementation of `handleFilePortListening`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



