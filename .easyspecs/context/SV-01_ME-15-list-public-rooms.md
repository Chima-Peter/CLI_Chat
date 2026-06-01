# Method ME-15 — ListPublicRooms

**Service:** SV-01 · **File:** SV-01_ME-15-list-public-rooms.md

## Summary
Lists all publicly available rooms on the server.

## Operation
`ListPublicRooms(cl *client)`
Iterates through all rooms hosted by the server (`s.rooms`), identifies those that are not private (`room_data.is_private == false`), and sends the list of public rooms to the client.

## Request / inputs
- `cl *client`: The client requesting the list of public rooms.

## Response / outputs
- If public rooms exist: Sends a message to the client with a JSON payload `{"rooms": [...]}` containing `room_id` and `room` name for each public room, along with a formatted string listing the room names.
- If no public rooms exist: Sends a message indicating "No publicly available rooms".

## Auth and permissions
No authentication or special permissions are required. The method checks the `is_private` flag of each room.

## Idempotency and concurrency
- The method uses `s.mu.RLock()` and `s.mu.RUnlock()` to ensure thread-safe access to the server's `rooms` map while iterating through it.
- This operation is read-only regarding the server state, making it inherently safe for concurrent calls.

## Errors
- No explicit error conditions for this method, as it simply returns a list or a "not found" message.

## Implementation notes
- The method is defined in `server/server.go`.
- It relies on the `server` struct's `rooms` map and `sync.RWMutex`.
- It uses `cl.send_user_message` to communicate with the client.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:263-288`: Implementation of `ListPublicRooms`.
- `server/server.go:19-24`: `server` struct definition including `rooms` map and `mu` (RWMutex).
- `server/rooms.go:10-20`: `room` struct definition, including `is_private` field.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



