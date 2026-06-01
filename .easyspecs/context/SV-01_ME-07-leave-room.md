# Method ME-07 — LeaveRoom

**Service:** SV-01 · **File:** SV-01_ME-07-leave-room.md

## Summary

The `LeaveRoom` method allows a client to leave a specific room, provided they are currently a member of that room. It handles member removal from the room structure, notifies the leaving client, and broadcasts the departure to other room members.

## Operation

The `LeaveRoom` operation is defined as a method on the server structure (`server/server.go:127-135`).

1. It resolves the room using `s.resolveRoom(roomID, roomName)`.
2. If resolution fails, it returns an error to the client.
3. If successful, it calls `room_data.LeaveRoom(cl)` to perform the departure logic (defined in `server/rooms.go:204-214`).

The `room.LeaveRoom` operation:
1. Verifies the client is part of the room (`server/rooms.go:205-208`).
2. Removes the member from the room via `r.DeleteMember(cl)` (`server/rooms.go:210`).
3. Notifies the client of the successful departure (`server/rooms.go:212`).
4. Broadcasts a message to the room announcing the departure (`server/rooms.go:213`).

## Request / inputs

- **`cl`**: `*client` - The client instance requesting to leave.
- **`roomID`**: `string` - The identifier of the room to leave.
- **`roomName`**: `string` - The name of the room to leave.

## Response / outputs

- **To the requester**: A message `"Left room."` via `cl.send_user_message`.
- **To the room members**: A broadcast message formatted as `"{client.nick} left the room"`.

## Auth and permissions

The method checks if the client is currently a member of the room (`server/rooms.go:205`). If not, an error is returned.

## Idempotency and concurrency

Not explicitly designed for idempotency. Attempting to leave a room when not a member will result in an error (`server/rooms.go:206-208`).

## Errors

- If the room cannot be resolved (`server/server.go:128-132`).
- If the client is not part of the room (`server/rooms.go:206-208`).

## Implementation notes

The implementation relies on `server.resolveRoom` to identify the target room and `room.DeleteMember` to perform the actual removal from the room's member list.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:127-135`: `LeaveRoom` method in `server` struct.
- `server/rooms.go:204-214`: `LeaveRoom` logic in `room` struct.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



