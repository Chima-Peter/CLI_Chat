# Method ME-25 — SendFriendRequest

**Service:** SV-01 · **File:** SV-01_ME-25-send-friend-request.md

## Summary

The `SendFriendRequest` method allows a client to send a friend request to another user, identified by their ID or name.

## Operation

The method is defined in the `server` struct, and delegates the request to the `client` object. It resolves the target user using `s.resolveUser(targetID, targetName)` and then calls `cl.SendFriendRequest(target)`.

Definition of `SendFriendRequest` in `server/server.go`:
```go
func (s *server) SendFriendRequest(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.SendFriendRequest(target)
}
```

Implementation logic in `server/client.go`:
```go
func (cl *client) SendFriendRequest(friend *client) {
    // ... validation logic (self, blocked, already friends, already sent, already pending) ...
    // ... state updates ...
    // ... messages to sender and recipient ...
}
```

## Request / inputs

- `cl`: The `*client` instance representing the user sending the request.
- `targetID`: The unique ID of the target user to send the friend request to.
- `targetName`: The name (nick) of the target user.

## Response / outputs

- The method does not return a direct value.
- It sends a JSON response to the calling client (`cl`) via `cl.send_user_message`:
    - On success: A `DONE` action message indicating the request was sent.
    - On failure (e.g., target not found, self-request, blocked, already friends): An `ERR` action message with an error description.
- It also sends a `DONE` action message to the target client if the request is successfully sent.

## Auth and permissions

- The sender must be authenticated.
- Validation checks:
    - Target must exist (resolved via `s.resolveUser`).
    - Sender cannot send a request to themselves.
    - Neither party can have blocked the other.
    - Sender and target cannot already be friends.
    - Sender cannot have already sent a request to target.
    - Sender cannot have a pending request from target.

## Idempotency and concurrency

- The method uses `cl.mu` and `friend.mu` (RWMutex) to ensure thread-safe access to the client's state (friends list, sent/pending requests, blocked users).
- The operation is generally idempotent in terms of the final state, but repeated requests before the previous one is processed/cleared will be rejected with an error message because of the "already sent" check.

## Errors

- "You cannot send a friend request to yourself."
- "You have blocked %s. Unblock them to send a friend request."
- "You cannot send a friend request to %s." (If the target has blocked the sender)
- "You are already friends with %s."
- "You have already sent a friend request to %s."
- "%s has already sent you a friend request."

## Implementation notes

The target user resolution is handled by `s.resolveUser`, which is a helper function in `server/utils.go`. The actual state management for friend requests involves updating `sent_friend_request` on the sender's client and `pending_friend_requests` on the recipient's client.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:576-583`: Definition of `SendFriendRequest` method in `server` struct.
- `server/client.go:69-113`: Implementation of `SendFriendRequest` in `client` struct, including validation and state updates.
- `server/utils.go:152`: Definition of `resolveUser` helper function.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



