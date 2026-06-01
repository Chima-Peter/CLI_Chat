# Method ME-04 — JoinRoom

**Service:** SV-01 · **File:** SV-01_ME-04-join-room.md

## Summary

The `JoinRoom` method allows a client to join a specified room by its ID or name. It handles room resolution and delegates the joining logic to the room instance.

## Operation

`func (s *server) JoinRoom(cl *client, roomID, roomName string)`

The method resolves the room using `s.resolveRoom(roomID, roomName)`. If successful, it invokes `room_data.JoinRoom(cl)` to perform the actual join operation, which includes validation of room membership, capacity, and privacy (password requirements).

## Request / inputs

- `cl`: The client attempting to join the room.
- `roomID`: The ID of the room to join (optional if `roomName` is provided).
- `roomName`: The name of the room to join (optional if `roomID` is provided).

## Response / outputs

- If successful, the client is added to the room's member list and receives a confirmation message.
- If the room is private, the client is prompted for a password (`GET_ROOM_PASSWORD`).
- If an error occurs (e.g., room not found, room full), the client receives an error message via `cl.err()`.

## Auth and permissions

- Any client can attempt to join a public room.
- If the room is private, a password is required (handled by the room's `JoinRoom` logic).

## Idempotency and concurrency

- The operation is generally idempotent regarding the final state (the client is in the room).
- Concurrency is managed using `sync.RWMutex` within the `server` and `room` structures to ensure thread-safe access to room data.

## Errors

- Returns an error if the room cannot be resolved (`resolveRoom`).
- The client receives an error if they are already in the room.
- The client receives an error if the room is full.

## Implementation notes

- Uses `s.resolveRoom` helper to locate the room data.
- Relies on `room.JoinRoom` for business logic related to joining a room instance.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:100-108`: Main JoinRoom handler definition.
- `server/utils.go:211-216`: Definition of `resolveRoom` helper.
- `server/rooms.go:76-115`: Definition of `room.JoinRoom` joining logic.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



