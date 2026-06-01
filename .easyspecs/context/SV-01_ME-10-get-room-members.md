# Method GetRoomMembers — Get room members

**Service:** SV-01 · **File:** SV-01_ME-10-get-room-members.md

## Summary

The `GetRoomMembers` method retrieves the list of members for a specified room and sends this information back to the requesting client.

## Operation

`GetRoomMembers(cl *client, roomID, roomName string)`

1. Resolves the room using `s.resolveRoom(roomID, roomName)`. If resolution fails, the error is sent to the client.
2. Checks if the requesting client (`cl`) is a member of the resolved room (`room_data.members`).
3. If the client is not a member, sends an error message: "Only members of a room can see it's members".
4. If the client is a member, fetches the list of room members using `room_data.FetchRoomMembers()`.
5. Sends the list of members back to the client via `cl.send_user_message`.

## Request / inputs

- `cl`: The client requesting the action.
- `roomID`: The ID of the room.
- `roomName`: The name of the room.

## Response / outputs

- The method sends a message to the client (`cl.send_user_message`) with the room ID, room name, and a formatted string containing the list of members.
- If an error occurs (e.g., room not found, client not a member), an error message is sent to the client (`cl.err`).

## Auth and permissions

- Requires the requesting client to be a member of the specified room.

## Idempotency and concurrency

- Read operation; idempotent.

## Errors

- Sends an error if the room cannot be resolved (`s.resolveRoom`).
- Sends an error if the client is not a member of the room.

## Implementation notes

- Implementation is in `server/server.go`.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:180-199`: `GetRoomMembers` method implementation.
- `server/server.go:181`: Room resolution.
- `server/server.go:187-191`: Membership check.
- `server/server.go:193`: Fetching members.
- `server/server.go:195-198`: Sending response.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



