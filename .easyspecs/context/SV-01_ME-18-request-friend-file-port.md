# Method ME-18 — requestFriendFilePort

**Service:** SV-01 · **File:** SV-01_ME-18-request-friend-file-port.md

## Summary

This method handles a request from a client to initiate a file transfer to a friend by requesting the friend's file listening port.

## Operation

`requestFriendFilePort(cl *client, friendName, filePath string)` in `server/server.go:403-441`.

## Request / inputs

- **`cl`**: The client requesting the file transfer.
- **`friendName`**: The name of the friend to whom the file is intended to be sent.
- **`filePath`**: The path of the file to be sent.

## Response / outputs

This method has no direct return value. It communicates results to the client:
- On validation or resolution errors, it calls `cl.err(error)`.
- If the friend's file port is not yet available, it queues the request and informs the client using `cl.send_user_message`.
- If the friend's file port is available, it calls `s.sendFilePortToSender`.

## Auth and permissions

- Requires the `friendName` and `filePath` to be provided and non-empty.
- The `friendName` must resolve to a valid existing user.
- The requesting client must be in the target friend's friend list (`hasID(friend.friends, cl.id)`).

## Idempotency and concurrency

- Uses `friend.mu.RLock()` / `friend.mu.RUnlock()` to safely check if the requester is in the friend's friend list.
- Uses `s.mu.Lock()` / `s.mu.Unlock()` when updating `s.pendingFileRequests` for queuing the request if the target friend's port is not immediately available.

## Errors

- "friend name is required" if `friendName` is empty or only whitespace.
- "file path is required" if `filePath` is empty or only whitespace.
- Errors returned by `s.resolveUser` if the friend cannot be found.
- "you are not in [friend's nick]'s friend list" if the friendship requirement is not met.

## Implementation notes

- Validates `friendName` and `filePath`.
- Resolves the friend's user object.
- Verifies that the requester is a friend of the target.
- Checks if the target friend is currently listening for file transfers (`friend.GetUserFilePort()`).
- If not listening, adds the request to `s.pendingFileRequests` and notifies the requester that the system is waiting for the friend.
- If listening, initiates the file port transfer via `s.sendFilePortToSender`.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:403-441`: Implementation of `requestFriendFilePort`.
- `server/server.go:414`: Resolves friend by name.
- `server/server.go:420-426`: Verifies friendship status.
- `server/server.go:430-435`: Queues pending file requests if port not available.
- `server/server.go:440`: Initiates file port transfer if available.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



