# Method ME-28 — CancelFriendRequest

**Service:** SV-01 · **File:** SV-01_ME-28-cancel-friend-request.md

## Summary

Cancels a previously sent friend request to another user.

## Operation

The `CancelFriendRequest` method is triggered on the server to remove a pending friend request that the caller (`cl`) previously sent to a `target`.

It performs the following steps:
1. Resolves the target user using `targetID` and `targetName` via `s.resolveUser`.
2. Validates the existence of the request.
3. Removes the request from the sender's sent requests list and the receiver's pending requests list.
4. Notifies both the sender and the receiver about the cancellation.

## Request / inputs

- **cl**: `*client` (The caller, implicitly from the context of the connection)
- **targetID**: `string` (The ID of the target user)
- **targetName**: `string` (The nick of the target user)

## Response / outputs

- None (void function, communicates results via messages to clients)

## Auth and permissions

- Implicitly authenticated as the caller.
- Requires that a friend request was actually sent to the target user.

## Idempotency and concurrency

- The operation uses `sync.RWMutex` to protect `cl.sent_friend_request` and `friend.pending_friend_requests` during read/write operations to ensure thread safety.

## Errors

- Errors related to user resolution (`s.resolveUser`) are sent to the client.
- Errors related to the absence of a sent friend request are sent to the client.

## Implementation notes

- The implementation in `server/server.go` orchestrates the resolution of the target user and then delegates the request deletion to `client.DeleteFriendRequest` in `server/client.go`.
- `client.DeleteFriendRequest` (lines 138-165 of `server/client.go`) handles the state update and notifications.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:603-610`: Implementation of `CancelFriendRequest` method in `server` struct.
- `server/client.go:138-165`: Implementation of `DeleteFriendRequest` method in `client` struct, which performs the actual state cleanup and notifications.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



