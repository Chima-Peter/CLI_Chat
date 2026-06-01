# Method ME-24 — ListMyRoomInvites

**Service:** SV-01 · **File:** SV-01_ME-24-list-my-room-invites.md

## Summary

Retrieves and lists all pending room invitations for the current user.

## Operation

`ListMyRoomInvites(cl *client)` (defined in `server/server.go:560`).

The method locks the client mutex to safely read `cl.room_invites`, iterates over them to collect room names, and then sends a formatted list of these names back to the client using `cl.send_user_message`.

## Request / inputs

- `cl *client`: The client instance invoking the method. The implementation assumes the client is authenticated and has a valid `room_invites` map.

## Response / outputs

- Sends a user message (`DONE` status) to the client:
    - If the user has no invites: "You have no room invites"
    - If the user has invites: A numbered list of the invites.

## Auth and permissions

- Assumes the client session is established.

## Idempotency and concurrency

- The method uses `cl.mu.RLock()` and `cl.mu.RUnlock()` to safely access `cl.room_invites` concurrently.
- It is a read-only operation and is idempotent.

## Errors

- No explicit error handling logic (e.g., if the user has no invites, it is handled as a message, not an error).

## Implementation notes

- Uses `cl.mu.RLock()` for concurrency safety.
- Relies on `cl.room_invites` mapping for data retrieval.
- Uses `formatNumberedList` helper for output formatting.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:560-574`: Implementation of ListMyRoomInvites.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



