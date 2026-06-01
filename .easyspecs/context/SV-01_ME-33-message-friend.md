# Method ME-33 — MessageFriend

**Service:** SV-01 · **File:** SV-01_ME-33-message-friend.md

## Summary

Sends a private message to a friend.

## Operation

`server.MessageFriend(cl *client, targetID, targetName, message string)`

This handler validates the target user, ensures the message content is not empty, verifies the friendship status, checks for blocking, and verifies the online status of the recipient before delivering the message.

## Request / inputs

- `cl`: The `*client` instance initiating the request.
- `targetID`: The unique identifier of the target friend.
- `targetName`: The nickname of the target friend.
- `message`: The message string to send.

## Response / outputs

The method delivers a JSON message to both the sender and the recipient via their respective `send_user_message` channels.
- To the sender: `DONE` action with a confirmation message.
- To the recipient: `MESSAGE_FRIEND` action containing the sender's details and the message content.

## Auth and permissions

- Requires the sender and recipient to be friends.
- Requires both users to not have blocked each other.
- The recipient must be online.

## Idempotency and concurrency

The method uses mutexes (`mu`) on both the sender's (`cl`) and the recipient's (`friend`) client structures to ensure thread-safe checks and state modifications, particularly for friendship and blocking status lookups.

## Errors

- Errors are reported back to the sender via `cl.err(err)` if:
    - User cannot be resolved (`targetID`/`targetName` mismatch or user not found).
    - Message content is empty.
    - Sender is not friends with the target.
    - Either user has blocked the other.
    - The recipient is offline.

## Implementation notes

- The `server` struct's `MessageFriend` method acts as the entry point, resolving the target user and then delegating the business logic to the `client`'s `MessageFriend` method.
- The `client`'s `MessageFriend` method handles the complex validation and message dispatching.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:642-649`: Implementation of the handler in the `server` struct.
- `server/client.go:265-321`: Implementation of the business logic in the `client` struct.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



