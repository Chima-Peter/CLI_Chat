# Method ME-03 — EditRoom

**Service:** SV-01 · **File:** SV-01_ME-03-edit-room.md

## Summary

Updates the configuration of an existing room, specifically its name and maximum allowed size.

## Operation

`EditRoom(cl *client, roomID, roomName, new_name string, max_size *int)`

This method first resolves the target room using `resolveRoom`. Then it delegates to the room's own `EditRoom` method to perform validation and state updates.

## Request / inputs

- `cl`: The client attempting the operation.
- `roomID`: Unique identifier for the room.
- `roomName`: Name of the room (alternative identifier).
- `new_name`: The new name to set for the room.
- `max_size`: A pointer to the new maximum size of the room.

## Response / outputs

- On success: A message `"Room settings updated."` with a payload containing the updated `room_id`, `room` (name), and `max_size`.
- On failure: An error message is returned to the client if the client is not the room owner, the new name is taken, or the new max size is invalid (negative or smaller than current members).

## Auth and permissions

Only the owner of the room (`r.owner.id == cl.id`) is authorized to edit the room.

## Idempotency and concurrency

The operation uses a mutex (`r.mu`) on the room object to protect against concurrent modification of the room's name and `max_size` fields.

## Errors

- `Only the room owner can edit this room.`
- `Max size cannot be negative.`
- `Cannot set max size below current member count (%d).`
- `Room name %s is already in use.`

## Implementation notes

The handler is defined in `server/server.go`. The actual business logic is in `server/rooms.go`.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:91-98`: Server handler method `EditRoom`.
- `server/rooms.go:157-202`: Implementation of `room.EditRoom`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



