# Method ME-34 — BlockUser

**Service:** SV-01 · **File:** SV-01_ME-34-block-user.md

## Summary

Blocks a user, preventing them from interacting with the current user.

## Operation

`server.BlockUser(cl *client, targetID, targetName string)`

This method is called on the server to block a specified target user. It resolves the target user's identity and then delegates the blocking action to the client instance of the caller.

## Request / inputs

- `cl`: The `client` instance initiating the block.
- `targetID`: The identifier of the user to be blocked.
- `targetName`: The nickname of the user to be blocked.

## Response / outputs

- If successful, sends a success message to the initiating client indicating that the user has been blocked.
- If unsuccessful, sends an error message to the initiating client.

## Auth and permissions

- Requires the initiating client to be authenticated (assumed context).
- Cannot block oneself.

## Idempotency and concurrency

- The operation checks if the user is already blocked in the `client.blocked_users` map under a read lock, and if not, adds them under a write lock, ensuring thread-safe access to the client's blocked user list.

## Errors

- `You cannot block yourself.`: If `targetID` matches the caller's ID.
- `%s is already blocked.`: If the target is already in the `blocked_users` map.
- Errors from `s.resolveUser` (e.g., if user is not found).

## Implementation notes

The server-side handler resolves the user first and then delegates to the client's `BlockUser` method. The client-side method handles the storage of the blocked user in a map.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:651-658`: Implementation of `server.BlockUser` handler.
- `server/client.go:323-342`: Implementation of `client.BlockUser` method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



