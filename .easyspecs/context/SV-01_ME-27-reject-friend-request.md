# Method ME-27 — RejectFriendRequest

**Service:** SV-01 · **File:** SV-01_ME-27-reject-friend-request.md

## Summary

The `RejectFriendRequest` method allows a user to decline a pending friend request they have received. It validates the existence of the request, removes it from both the recipient's pending list and the sender's sent list, and notifies both parties of the outcome.

## Operation

The operation is handled by `server.RejectFriendRequest` which resolves the target user and then invokes the client-side logic `client.RejectFriendRequest`.

## Request / inputs

- **`cl`**: `*client` (the user performing the rejection).
- **`targetID`**: `string` (the ID of the user whose request is being rejected).
- **`targetName`**: `string` (the nickname of the user whose request is being rejected).

## Response / outputs

- **Success:** Both the recipient (`cl`) and the sender (`friend`) receive a notification message:
  - Recipient: `"Friend request from %s rejected."`
  - Sender: `"%s rejected your friend request."`
- **Error:** If the request does not exist, an error is returned to the recipient.

## Auth and permissions

- Requires an active session for the calling client.
- The user must have a pending friend request from the specified `targetID`/`targetName`.

## Idempotency and concurrency

- The operation uses mutexes (`cl.mu`, `friend.mu`) to safely modify the friend request lists, ensuring concurrency safety during the deletion of the request.
- The check for existence (`hasID`) before deletion makes the operation safe to perform; if the request is not present, it gracefully returns an error.

## Errors

- **User resolution failure:** `cl.err(err)` if `resolveUser` fails.
- **Request not found (recipient):** `cl.err(fmt.Errorf("You have not received a friend request from this user."))`
- **Request not found (sender integrity check):** `cl.err(fmt.Errorf("Error rejecting sent requests. Record not found."))`

## Implementation notes

- The `server` handles user resolution, while the `client` handles the state manipulation and notification.
- The lists of pending requests and sent requests are protected by RWMutexes.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:594-601`: Handler implementation for RejectFriendRequest.
- `server/client.go:167-194`: Core logic for RejectFriendRequest in client struct.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



