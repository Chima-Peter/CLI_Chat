# Method ME-20 — dispatchPendingFileRequests

**Service:** SV-01 · **File:** SV-01_ME-20-dispatch-pending-file-requests.md

## Summary

The `dispatchPendingFileRequests` method processes and dispatches queued file transfer requests for a specific client (`receiver`). It facilitates the setup of a file transfer connection by connecting a sender to the receiver's open file port.

## Operation

This is an internal helper method of the `server` struct (`server/server.go`). It:
1.  Atomically retrieves and removes pending file requests for the given `receiver` from the server's `pendingFileRequests` map.
2.  Obtains the receiver's file listening host and port.
3.  Iterates through all pending requests.
4.  Validates that the sender still exists and is a friend of the receiver.
5.  If validation passes, invokes `sendFilePortToSender` to communicate the receiver's listening host and port to the sender.

## Request / inputs

- `receiver *client`: The client object for whom pending file requests should be dispatched.

## Response / outputs

- None. The method performs actions via side effects (`sendFilePortToSender`).

## Auth and permissions

- Requires that the requester (sender of the file) is currently in the `friends` list of the receiver.

## Idempotency and concurrency

- Uses `s.mu.Lock()` to safely access and modify `s.pendingFileRequests`.
- Uses `receiver.mu.RLock()` to safely access `receiver.friends`.
- Uses `s.mu.RLock()` to safely access `s.clients`.

## Errors

- Silently skips requests if the sender does not exist in `s.clients` or is not a friend of the receiver.
- Terminates early if the receiver has no pending requests or if the receiver's file port is not configured (host empty or port 0).

## Implementation notes

- The implementation resides in `server/server.go`.
- Relies on `receiver.GetUserFilePort()` to determine where to connect for file transfer.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:470-501`: Implementation of `dispatchPendingFileRequests` method.
- `server/server.go:503-511`: Implementation of `sendFilePortToSender` invoked by `dispatchPendingFileRequests`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



