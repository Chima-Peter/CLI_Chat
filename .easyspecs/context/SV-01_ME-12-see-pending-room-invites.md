# Method ME-12 — SeePendingRoomInvites

**Service:** SV-01 · **File:** SV-01_ME-12-see-pending-room-invites.md

## Summary
Allows a room owner to retrieve and view the list of pending room invites for a specified room.

## Operation
`server.SeePendingRoomInvites(cl, roomID, roomName)` -> `room.SeePendingRoomInvites(owner)`

- Validates room resolution and admin ownership.
- Fetches and formats the list of invited users.

## Request / inputs
- `client` (contextual, implicit)
- `roomID` (string, room identifier)
- `roomName` (string, room name identifier)

## Response / outputs
- A formatted user message sent to the client:
  - `DONE` message type
  - Body: A numbered list of invited user nicknames.

## Auth and permissions
- Requires that the requesting client be the **owner** of the specified room (`room_data.owner.id == cl.id`).

## Idempotency and concurrency
- Uses `r.mu.RLock()` to safely access the room's `invites` map concurrently.

## Errors
- Returns error if the room cannot be resolved (`resolveRoom` error).
- Returns error if the client is not the room admin ("This action is reserved for only the admin!").

## Implementation notes
- The implementation is split between `server/server.go` (auth and delegation) and `server/rooms.go` (retrieval and formatting).

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:228-241`: Implementation of `SeePendingRoomInvites` handler in `server` struct, handling authorization and delegation.
- `server/rooms.go:291-300`: Implementation of `SeePendingRoomInvites` in `room` struct, handling concurrency and response formatting.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



