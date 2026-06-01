# Method ME-23 — ListMyRooms

**Service:** SV-01 · **File:** SV-01_ME-23-list-my-rooms.md

## Summary

Lists all the rooms that the current authenticated client is a member of.

## Operation

`ListMyRooms(cl *client)`

It iterates over the client's `my_rooms` field, constructs a list of maps containing the `room_id` and `room` (name) for each room, and sends this data back to the client as a `rooms` object, along with a formatted string message listing the room names.

## Request / inputs

- `cl *client`: The authenticated client requesting the list of their rooms.

## Response / outputs

- None (void method).
- Sends a response to the client (`cl.send_user_message`) with:
    - Payload: `{"rooms": []map[string]string}` where each map contains `room_id` and `room`.
    - Status/Message: `DONE` status with a formatted list of room names.

## Auth and permissions

- Requires an authenticated client context, provided as `*client`.

## Idempotency and concurrency

- The method does not modify any state, so it is inherently idempotent.
- It accesses `cl.my_rooms`. 

## Errors

- No explicit error handling.

## Implementation notes

- Uses `cl.send_user_message` to communicate the result back to the client.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:545-558`: Implementation of ListMyRooms method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



