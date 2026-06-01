# Method ME-35 — UnblockUser

**Service:** SV-01 · **File:** SV-01_ME-35-unblock-user.md

## Summary

The `UnblockUser` method allows a client to unblock another user, removing them from the caller's blocked users list.

## Operation

`server.(*server).UnblockUser(cl *client, targetID, targetName string)`

This method resolves the target user using `s.resolveUser(targetID, targetName)` and then calls `cl.UnblockUser(target)` on the client instance to perform the actual unblocking.

## Request / inputs

- `cl`: The client requesting the unblock action.
- `targetID`: The unique ID of the user to unblock.
- `targetName`: The nickname of the user to unblock.

## Response / outputs

- None (void). Any errors are sent to the client via `cl.err(err)`.

## Auth and permissions

- Requires a valid active connection (implicit as `cl` is provided).
- Checks if the user is actually in the blocked list.

## Idempotency and concurrency

- `cl.UnblockUser` uses a mutex (`cl.mu.Lock()`) to safely modify the `blocked_users` list, ensuring thread safety.

## Errors

- Errors occur if the user cannot be resolved (`resolveUser` fails).
- Errors occur if the user is not currently blocked (`cl.err` is called with "is not blocked.").

## Implementation notes

The method performs a two-step process: first, resolving the target user identity, then invoking the client-specific unblock logic which handles the internal blocked user list management.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:660-667`: Implementation of `UnblockUser` on the server.
- `server/client.go:344-351`: Implementation of `UnblockUser` on the client (business logic).

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



