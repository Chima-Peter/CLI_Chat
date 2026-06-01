# Method ME-29 — SeePendingFriendRequests

**Service:** SV-01 · **File:** SV-01_ME-29-see-pending-friend-requests.md

## Summary

The `SeePendingFriendRequests` method retrieves and sends the list of pending friend requests for a given client.

## Operation

The method `SeePendingFriendRequests` is defined on the `*server` struct:

```go
func (s *server) SeePendingFriendRequests(cl *client) {
	cl.mu.RLock()
	ids := copyIDSet(cl.pending_friend_requests)
	cl.mu.RUnlock()
	cl.FetchPendingRequests(s.usersFromIDs(ids))
}
```

## Request / inputs

- `cl`: A pointer to the `client` requesting their pending friend requests.

## Response / outputs

This method does not return a value. Instead, it populates the client's pending friend requests state by calling `cl.FetchPendingRequests` with a list of users corresponding to the pending request IDs.

## Auth and permissions

The method assumes the `client` is already authenticated and associated with a valid session context as it operates on the provided `cl` object directly.

## Idempotency and concurrency

- **Concurrency**: The method uses `cl.mu.RLock()` and `cl.mu.RUnlock()` to safely copy the set of pending friend request IDs from the client, ensuring thread-safe access to the client's state.
- **Idempotency**: The operation is a read-like operation (fetching state to the client) and is idempotent.

## Errors

There are no explicit error handling mechanisms for this method.

## Implementation notes

The method performs the following steps:
1. Acquires a read lock on the client (`cl.mu.RLock()`).
2. Copies the set of pending friend request IDs (`cl.pending_friend_requests`).
3. Releases the read lock (`cl.mu.RUnlock()`).
4. Retrieves user objects from IDs and triggers `cl.FetchPendingRequests`.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:612-617`: Implementation of the `SeePendingFriendRequests` method on the `server` struct.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



