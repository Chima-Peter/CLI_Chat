# Method ME-32 — DeleteFriend

**Service:** SV-01 · **File:** SV-01_ME-32-delete-friend.md

## Summary

The `DeleteFriend` method handles the removal of a friend from a user's friend list on the server.

## Operation

It takes a client and target user identifiers as input, resolves the target user, and then invokes the client's `DeleteFriend` method.

Defined in `server/server.go:633-640`.

## Request / inputs

- `cl *client`: The client initiating the action.
- `targetID string`: The ID of the friend to be deleted.
- `targetName string`: The username of the friend to be deleted.

## Response / outputs

This method performs the action on the server side and initiates further actions on the client. It does not return a direct value but uses `cl.err(err)` to report any errors during user resolution.

## Auth and permissions

The method assumes the caller is authenticated as the client `cl`. The `resolveUser` method implicitly handles checking the validity of the target user identifiers.

## Idempotency and concurrency

The operation relies on `resolveUser` and `cl.DeleteFriend`. Concurrency handling depends on the implementation of `cl.DeleteFriend`.

## Errors

If `s.resolveUser(targetID, targetName)` fails (e.g., target user not found), the error is reported back to the client using `cl.err(err)`.

## Implementation notes

The implementation directly calls `resolveUser` to identify the target user and then delegates the deletion to the client's own `DeleteFriend` method.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:633-640`: Implementation of the `DeleteFriend` method on the `server` struct.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



