# Method ME-16 — SwitchContext

**Service:** SV-01 · **File:** SV-01_ME-16-switch-context.md

## Summary

The `SwitchContext` method allows a client to switch their active communication context between a room or a friend. This dictates where subsequent messages or actions are directed.

## Operation

`server.SwitchContext(cl *client, contextType, name string)`

This method changes the current context of a client. It validates the existence of the target (room or user) and ensures the client has appropriate permissions (i.e., is a member of the room or friends with the user).

## Request / inputs

- `cl *client`: The client initiating the context switch.
- `contextType string`: The type of context to switch to ("room" or "friend").
- `name string`: The name of the room or friend to switch to.

## Response / outputs

- If successful, it sends a message to the client confirming the switch to the specified room or friend using `cl.send_user_message`.
- If unsuccessful (e.g., name is empty, room/user not found, not a member/friend), it sends an error message to the client using `cl.err`.

## Auth and permissions

- Validates that the client is indeed in the specified room or friends with the specified user.
- Uses `s.resolveRoom` or `s.resolveUser` to locate the target.

## Idempotency and concurrency

- Uses `cl.mu.Lock()` and `cl.mu.Unlock()` to synchronize access to the client's state (`current_context`, `room`, `current_friend`) during the switch.

## Errors

- Returns errors if:
    - `name` is empty.
    - `contextType` is invalid (not "room" or "friend").
    - Target room not found or client is not a member of the room.
    - Target user not found or client is not friends with the user.

## Implementation notes

The method performs input sanitation (lowercasing, trimming) on `contextType` and `name` before processing.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:301-367`: Implementation of the `SwitchContext` method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



