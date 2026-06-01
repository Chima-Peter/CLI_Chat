# Method ME-08 — DeleteRoom

**Service:** SV-01 · **File:** SV-01_ME-08-delete-room.md

## Summary

The `DeleteRoom` method removes an existing room from the server.

## Operation

The operation is performed by `DeleteRoom(cl *client, roomID, roomName string)` in `server/server.go`.
It resolves the room using `s.resolveRoom` and then calls `DeleteRoom` on the room data.

## Request / inputs

- `cl`: The client requesting the deletion.
- `roomID`: The identifier of the room to delete.
- `roomName`: The name of the room to delete.

## Response / outputs

No explicit return value. Errors, if any, are reported back to the client using `cl.err(err)`.

## Auth and permissions

Relies on `s.resolveRoom` (checks if room exists and user has access) and `room_data.DeleteRoom` (checks if user is allowed to delete the room).

## Idempotency and concurrency

`s.resolveRoom` is likely thread-safe (assuming `server` handles locking), and `room_data.DeleteRoom` should perform the deletion atomically.

## Errors

Errors from `resolveRoom` (e.g., room not found, permission denied) are passed to the client via `cl.err(err)`.

## Implementation notes

- The implementation in `server/server.go` relies on helper methods for room resolution and deletion.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:137-145`: Implementation of DeleteRoom method.
- `server/utils.go:211`: Implementation of resolveRoom helper method.
- `server/rooms.go:216`: Implementation of room.DeleteRoom method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



