# Method ME-30 — SeeSentFriendRequests

**Service:** SV-01 · **File:** SV-01_ME-30-see-sent-friend-requests.md

## Summary

Retrieves a list of friend requests sent by the current client.

## Operation

`SeeSentFriendRequests` on `server` struct:
```go
func (s *server) SeeSentFriendRequests(cl *client) {
	cl.mu.RLock()
	ids := copyIDSet(cl.sent_friend_request)
	cl.mu.RUnlock()
	cl.FetchSentRequests(s.usersFromIDs(ids))
}
```

## Request / inputs

- `cl *client`: The client requesting the list.

## Response / outputs

- The client's `FetchSentRequests` method is called with a slice of user objects corresponding to the IDs of the sent friend requests.

## Auth and permissions

No specific permissions checked, implicitly authorized by the authenticated client session.

## Idempotency and concurrency

Uses `cl.mu.RLock()` to safely access the client's `sent_friend_request` set.

## Errors

No specific error handling in this method.

## Implementation notes

Copies the set of IDs (`cl.sent_friend_request`), resolves them to user objects via `s.usersFromIDs(ids)`, and passes the resulting user list to the client's `FetchSentRequests` method.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:619-624`: `SeeSentFriendRequests` method implementation.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



