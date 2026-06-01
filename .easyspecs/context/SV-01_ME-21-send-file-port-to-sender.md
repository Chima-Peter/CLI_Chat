# Method ME-21 — sendFilePortToSender

**Service:** SV-01 · **File:** SV-01_ME-21-send-file-port-to-sender.md

## Summary

Sends a file port notification to the sender, indicating that the receiver is ready to receive a file at a specific host and port.

## Operation

`func (s *server) sendFilePortToSender(sender, receiver *client, host string, port int, filePath string)`

The method constructs a message containing the receiver's connection details (host, port, nickname) and the file path, and sends it to the sender client, notifying them that the file port is ready for a transfer.

## Request / inputs

- `sender` (*client): The client initiating the file transfer.
- `receiver` (*client): The client intending to receive the file.
- `host` (string): The host address where the receiver is listening for the file transfer.
- `port` (int): The port number where the receiver is listening.
- `filePath` (string): The path of the file to be transferred.

## Response / outputs

None. The operation has a side-effect: it sends a `FILE_PORT_LISTENING` message to the `sender` client.

## Auth and permissions

No explicit auth checks within this method; it relies on the caller (`requestFriendFilePort` or `dispatchPendingFileRequests`) ensuring the validity of the relationship between `sender` and `receiver`.

## Idempotency and concurrency

- **Idempotency**: N/A.
- **Concurrency**: Relies on the `sender` client's thread-safe message sending mechanism (`send_user_message`).

## Errors

No explicit error handling in this method.

## Implementation notes

This method is internal to the `server` struct. It is invoked when a file transfer request is initiated or fulfilled after a pending request is resolved.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:503-511`: Implementation of `sendFilePortToSender`.
- `server/server.go:440`: Usage in `requestFriendFilePort`.
- `server/server.go:499`: Usage in `dispatchPendingFileRequests`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



