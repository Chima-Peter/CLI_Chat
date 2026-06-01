# Method ME-01 — CreateRoom

**Service:** SV-01 · **File:** SV-01_ME-01-create-room.md

## Summary

The `CreateRoom` method allows an authenticated client to create a new chat room on the server. The room is initially public, and the creating client becomes the room's owner.

## Operation

`server.CreateRoom(cl *client, room_name string)`

This method is a handler on the server instance that processes a request to create a new room. It validates the room name, ensures it's unique, initializes a `room` struct, registers it with the server, adds the client to the room, and prompts the client to set a password (optional).

## Request / inputs

- `cl *client`: The client requesting the room creation.
- `room_name string`: The desired name for the new room.

## Response / outputs

There is no direct return value. The method performs side effects:

- Updates the server's room list (`s.rooms`).
- Updates the client's room list (`cl.my_rooms`).
- Sets the client's current room (`cl.room`).
- Sends a `SET_ROOM_PASSWORD` message to the client, providing the new `room_id` and the room name.

## Auth and permissions

- Assumes the client is already authenticated.
- No explicit permission check other than that the client must be valid.

## Idempotency and concurrency

- The method uses `sync.RWMutex` to manage access to the `server.rooms` map, ensuring safe concurrent access.
- It is not strictly idempotent as it creates a new room with a unique UUID every time it is called.

## Errors

- Returns an error to the client if:
  - The room name is empty (after trimming).
  - The room name is already taken.

## Implementation notes

- Uses `strings.TrimSpace` to sanitize the room name.
- Generates a UUID for the new room.
- Defaults `is_private` to `false`.
- The client who creates the room is added as the owner and automatically joins the room.
- Sends a `SET_ROOM_PASSWORD` instruction as the final step.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:36-74`: CreateRoom method implementation.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



