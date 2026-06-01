# Method ME-17 — SendContextMessage

**Service:** SV-01 · **File:** SV-01_ME-17-send-context-message.md

## Summary

The `SendContextMessage` method facilitates sending messages within the user's current interaction context, which can be either a room or a direct conversation with a friend.

## Operation

The method determines the target (room or friend) based on the client's current context and delegates the message handling accordingly.

- `server/server.go:369-401`: `SendContextMessage` definition.

## Request / inputs

- `cl *client`: The client session initiating the message.
- `message string`: The message text to send.

## Response / outputs

- None (void).
- Communicates via `cl.err` for errors and `cl.send_user_message` or `cl.MessageFriend` for successful operations.

## Auth and permissions

The method assumes the client session is already established and authenticated. It relies on the internal `cl.current_context`, `cl.room`, and `cl.current_friend` state variables, which are set via other operations (presumably `SwitchContext`).

## Idempotency and concurrency

- Uses `cl.mu.RLock()` / `cl.mu.RUnlock()` to read the client's context and related state (`cl.room`, `cl.current_friend`) safely in a concurrent environment.

## Errors

- "message cannot be empty": If the input `message` is empty or whitespace only.
- "use /switch room <name> before sending a message": If context is `contextRoom` but `cl.room` is `nil`.
- "use /switch friend <name> before sending a message": If context is `contextFriend` but `cl.current_friend` is `nil`.
- "use /switch room|friend <name> before sending a message": If the current context is not recognized.

## Implementation notes

The method performs a `strings.TrimSpace` on the message before validation.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:369-401`: Implementation of `SendContextMessage`.
- `server/server.go:376-380`: Concurrent access to client state with mutex lock.
- `server/server.go:382-400`: Context-based dispatching logic.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



