# Method ME-02 — SetRoomPassword

**Service:** SV-01 · **File:** SV-01_ME-02-set-room-password.md

## Summary

The `SetRoomPassword` method allows the owner of a room to update the password required for joining that room.

## Operation

`server.SetRoomPassword(cl *client, roomID, roomName, password string)` updates the room's password if the authenticated client is the owner of the room.

## Request / inputs

- `cl`: The client requesting the password update.
- `roomID`: The unique identifier of the room (used to resolve the room).
- `roomName`: The name of the room (used to resolve the room).
- `password`: The new password to set for the room.

## Response / outputs

This method does not return a value. Instead, it interacts directly with the provided `client` instance to report errors or progress (via `cl.err`).

## Auth and permissions

The method checks if the requesting client is the owner of the room by comparing `room_data.owner.id` with `cl.id`. If they do not match, the action is denied.

## Idempotency and concurrency

The implementation relies on `s.resolveRoom` and `room_data.SetRoomPassword`. Its idempotency depends on the underlying implementation of `room_data.SetRoomPassword`.

## Errors

Errors are communicated back to the client via `cl.err(err)`.
- If the room cannot be resolved (`resolveRoom` returns an error), that error is propagated.
- If the client is not the owner, a specific error: "You don't have the right to update this room password." is returned.

## Implementation notes

The method performs a room resolution, an ownership check, and then calls `room_data.SetRoomPassword` to perform the update.

## Revision

- Initial draft: contract and handler from implementation.
- Added resolveRoom helper reference to evidence index.

## Evidence index

- `server/server.go:76-89`: Implementation of `SetRoomPassword` method.
- `server/utils.go:211`: Implementation of `resolveRoom` helper method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



