# Method ME-06 — JoinRoomWithPassword

**Service:** SV-01 · **File:** SV-01_ME-06-join-room-with-password.md

## Summary

The `JoinRoomWithPassword` method allows a client to join a private room by providing the correct password.

## Operation

`server.JoinRoomWithPassword(cl *client, roomID, roomName, password string)`

This method processes a request to join a room protected by a password. It validates the password input, resolves the target room, and delegates the joining process to the room instance if the password is valid.

## Request / inputs

- `cl`: The `client` object initiating the request.
- `roomID`: The unique identifier of the room to join (optional if `roomName` is provided).
- `roomName`: The name of the room to join (optional if `roomID` is provided).
- `password`: The password to join the room.

## Response / outputs

- This method does not return a value directly. It interacts with the client via the `cl.err(error)` method if the password is invalid, the room cannot be resolved, or the join operation fails.

## Auth and permissions

- Requires a valid password to join the room.
- Access control is enforced by the room's `JoinRoomWithPassword` implementation.

## Idempotency and concurrency

- The idempotency and concurrency management depends on the implementation of `room.JoinRoomWithPassword`.

## Errors

- Returns an error to the client if the provided password is empty after trimming.
- Returns an error if the room cannot be resolved using `s.resolveRoom(roomID, roomName)`.
- May return errors from the room's `JoinRoomWithPassword` method (e.g., incorrect password, full room).

## Implementation notes

- Password is trimmed of whitespace before validation.
- Room resolution is handled by `s.resolveRoom(roomID, roomName)`.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:110-125`: Implementation of the `JoinRoomWithPassword` method.
- `server/utils.go:211-230`: Implementation of the `resolveRoom` helper method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



