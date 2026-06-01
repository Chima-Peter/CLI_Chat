# Method ME-09 — DeleteRoomMember

**Service:** SV-01 · **File:** SV-01_ME-09-delete-room-member.md

## Summary
Deletes a member from a specified room. This operation is restricted to the room owner.

## Operation
`DeleteRoomMember(cl *client, roomID, roomName, memberID, memberName string)`
Implemented in `server/server.go:147-178`.

## Request / inputs
- `cl`: The client requesting the action.
- `roomID`: ID of the room.
- `roomName`: Name of the room.
- `memberID`: ID of the member to be removed.
- `memberName`: Nickname of the member to be removed.

## Response / outputs
- Upon success, broadcasts a message to the room, sends a message to the removed member, and sends a confirmation to the caller (`cl`).
- Errors are returned to the caller (`cl`) via `cl.err(err)`.

## Auth and permissions
Restricted to the owner of the room (`room_data.owner.id == cl.id`).

## Idempotency and concurrency
Not explicitly designed for idempotency.

## Errors
- `s.resolveRoom` failure (e.g., room not found).
- Unauthorized (caller is not owner).
- `s.resolveUser` failure (e.g., user not found).
- Member not in the room.

## Implementation notes
Relies on `s.resolveRoom` and `s.resolveUser`. Uses `room_data.DeleteMember(member)`. Notifies participants via `member.send_user_message`, `room_data.Broadcast`, and `cl.send_user_message`.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:147-178`: Implementation of DeleteRoomMember.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



