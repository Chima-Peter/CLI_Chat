# Method ME-36 — GetUserStatus

**Service:** SV-01 · **File:** SV-01_ME-36-get-user-status.md

## Summary

The `GetUserStatus` method retrieves and returns the current status (online/offline), room connection details, and block status of a specified user.

## Operation

`server.GetUserStatus(cl *client, targetID, targetName string)`

This method first resolves the target user using `s.resolveUser(targetID, targetName)`. If successful, it invokes `cl.GetUserStatus(target)` to retrieve the target's status and send it to the requester.

## Request / inputs

- **cl**: The client requesting the status.
- **targetID**: The ID of the target user.
- **targetName**: The nickname of the target user.

## Response / outputs

The method does not return a direct value. Instead, it sends a message (with action `DONE`) to the requesting client (`cl`) containing the following payload:
- `user_id`: Target user ID.
- `nick`: Target user nickname.
- `online`: Boolean indicating online status.
- `blocked`: Boolean indicating if the target is blocked by the requester.
- `room_id`: The ID of the room the target is currently in (or empty string).
- `room`: The name of the room the target is currently in (or empty string).

## Auth and permissions

No specific authorization is required to check a user's status, but the target user must exist.

## Idempotency and concurrency

The operation is read-only and therefore idempotent. It uses `RWMutex` to ensure safe access to shared state (such as the client's blocked users map and the target's online status and room information).

## Errors

- Returns an error if the target user cannot be resolved. The error is sent to the requesting client via `cl.err(err)`.

## Implementation notes

The method relies on `resolveUser` to identify the target client and then uses `client.GetUserStatus` to gather information from the target's state.

## Revision

- Initial draft: contract and handler from implementation.
- Fixed Evidence index to remove forbidden citations.

## Evidence index

- `server/server.go:669-676`: `server.GetUserStatus` method definition.
- `server/client.go:366-394`: `client.GetUserStatus` implementation.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



