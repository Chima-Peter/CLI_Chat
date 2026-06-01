# Field FD-11 — sent_friend_request

**Entity:** DM-02 · **File:** DM-02_FD-11-sent-friend-request.md

## Summary

The `sent_friend_request` field maintains a set of user identifiers (as a map to empty structures) representing the friend requests that the client has initiated but that have not yet been accepted or rejected.

## Type and constraints

- **Type:** `map[string]struct{}`
- **Constraints:** The keys in the map are the unique string identifiers (`id`) of the target users. The value `struct{}` is used to effectively implement a set for O(1) membership lookups.
- **Concurrency:** Access to this field is protected by the client's `sync.RWMutex` (`mu`). Read operations must use `RLock()` and write operations must use `Lock()`.

## Default and nullability

- **Default:** An empty map upon client initialization.
- **Nullability:** The map itself is initialized as a non-nil, empty map; it is not typically set to `nil`.

## Validation

- The system ensures that a client cannot send a friend request to themselves (`server/client.go:70`).
- The system verifies that the sender has not already blocked the target user, and vice versa, before allowing a request (`server/client.go:75-83`).
- The system checks if the users are already friends, if a request was already sent, or if a pending request from the target already exists, to prevent duplicate or redundant requests (`server/client.go:86-100`).

## Privacy / sensitivity

- This field contains information about the social graph of the client (users they have interacted with by initiating a friend request). It should be treated as sensitive user data.

## Representation in API and UI

- The `FetchSentRequests` method retrieves the list of users to whom a friend request has been sent and sends this information back to the client via `send_user_message` (`server/client.go:196-208`).

## Revision

- <!-- Initial draft: field mapping from server/client.go -->

## Evidence index

- `server/client.go:28`: Field definition `sent_friend_request map[string]struct{}`.
- `server/client.go:104`: Usage of `sent_friend_request` for storing new outgoing friend requests.
- `server/client.go:140`: Usage of `sent_friend_request` for checking existence of outgoing friend requests in `DeleteFriendRequest`.
- `server/client.go:196-208`: Usage of `sent_friend_request` (implicitly via logic using the field) in `FetchSentRequests`.
- `server/client.go:34`: Mutex `mu` protecting the field.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



