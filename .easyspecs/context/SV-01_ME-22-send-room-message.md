# Method ME-22 — SendRoomMessage

**Service:** SV-01 · **File:** SV-01_ME-22-send-room-message.md

## Summary

The `SendRoomMessage` method allows a client to broadcast a message to a specific room they are currently a member of.

## Operation

`func (s *server) SendRoomMessage(cl *client, roomName, message string)` defined in `server/server.go:513-543`.

## Request / inputs

- `cl`: The `*client` object representing the user initiating the message.
- `roomName`: The name of the room to which the message should be sent.
- `message`: The content of the message to broadcast.

## Response / outputs

- If successful, it broadcasts the message to the room via `room_data.Broadcast(cl, message)` and sends a confirmation back to the client (`cl.send_user_message`).
- Returns errors directly to the client via `cl.err()` if the message is empty, the room name is missing, or if the client is not in the specified room.

## Auth and permissions

- Validates that the client is a member of the requested room by checking `cl.my_rooms[room_data.id]`.

## Idempotency and concurrency

- Not idempotent; sending a message triggers a broadcast to all members in the room.

## Errors

- "message cannot be empty": Raised if the provided message is empty or only whitespace.
- "room name is required": Raised if the room name is empty.
- "you are not in room: %s": Raised if the client is not a member of the room or the room cannot be resolved.

## Implementation notes

- The method trims whitespace from the message.
- It attempts to find the room using `cl.roomByName`.
- If not found in the client's cached rooms, it performs a lookup using `s.resolveRoom` and subsequently checks for membership in `cl.my_rooms`.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:513-543`: Implementation of the SendRoomMessage method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



