# Method ME-26 — AcceptFriendRequest

**Service:** SV-01 · **File:** SV-01_ME-26-accept-friend-request.md

## Summary

Accepts a pending friend request from another user. This method removes the request from the pending list for the receiver and the sent list for the sender, then updates both users' friend lists and notifies them of the success.

## Operation

The method `AcceptFriendRequest` is defined in `server/server.go`. It takes a `client` and target user identification (`targetID`, `targetName`) as arguments. It resolves the target user and then invokes `client.AcceptFriendRequest` on the receiver's client object.

## Request / inputs

- `targetID` (string): The ID of the user whose friend request is to be accepted.
- `targetName` (string): The name of the user whose friend request is to be accepted.

## Response / outputs

- If successful, both users are added to each other's friend lists, the pending request is removed, and both users receive a confirmation message.
- If the friend request does not exist, an error is returned to the client.

## Auth and permissions

The operation assumes the requesting client is already authenticated and authorized to perform actions within their session.

## Idempotency and concurrency

The method is idempotent in the sense that it first checks if the request exists (`hasID`) before attempting to accept it. It uses mutex locks (`cl.mu`, `friend.mu`) to ensure thread safety when modifying friend lists and request status.

## Implementation notes

The `server.AcceptFriendRequest` function resolves the user before calling `client.AcceptFriendRequest`.

```go
// server/server.go:585-592
func (s *server) AcceptFriendRequest(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.AcceptFriendRequest(target)
}
```

The `client.AcceptFriendRequest` function handles the logic:

```go
// server/client.go:115-136
func (cl *client) AcceptFriendRequest(friend *client) {
	cl.mu.RLock()
	exists := hasID(cl.pending_friend_requests, friend.id)
	cl.mu.RUnlock()
	if !exists {
		cl.err(fmt.Errorf("You have not received a friend request from this user."))
		return
	}

	cl.mu.Lock()
	cl.friends[friend.id] = struct{}{}
	delete(cl.pending_friend_requests, friend.id)
	cl.mu.Unlock()

	friend.mu.Lock()
	friend.friends[cl.id] = struct{}{}
	delete(friend.sent_friend_request, cl.id)
	friend.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("You are now friends with %s.", friend.nick))
	friend.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s accepted your friend request.", cl.nick))
}
```

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:585-592`: Implementation of `server.AcceptFriendRequest`.
- `server/client.go:115-136`: Implementation of `client.AcceptFriendRequest`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



