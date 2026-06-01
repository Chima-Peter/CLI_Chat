# Field FD-10 — pending_friend_requests

**Entity:** DM-02 · **File:** DM-02_FD-10-pending-friend-requests.md

## Summary

The `pending_friend_requests` field maintains a collection of unique user identifiers (`client.id`) for clients that have sent a friend request to the current client but have not yet been accepted or rejected.

## Type and constraints

- **Type:** `map[string]struct{}`
- **Constraint:** The keys of the map are `string` identifiers corresponding to the `id` of other `client` objects. The values are empty `struct{}` types, used efficiently as a set to check for existence.

## Default and nullability

- **Default:** An initialized empty map.
- **Nullability:** Not nullable; implemented as a Go map type, which can be safely accessed (read) when nil, but requires initialization before writing.

## Validation

- **Existence Checks:** Before performing actions such as accepting or rejecting a request, the system verifies the presence of the requester's ID in this map using a helper function (e.g., `hasID`).
- **Mutual Exclusivity:** Logic ensures that a client cannot have a pending request from someone they have already blocked, or from someone to whom they have already sent a request.

## Privacy / sensitivity

- **Sensitivity:** This field contains information about pending friend requests, which is considered private to the user.
- **Access Control:** Interactions with this field are mediated through the `client` struct methods, which are protected by the `client.mu` `sync.RWMutex` to ensure thread safety during concurrent operations.

## Representation in API and UI

- **API:** When fetching pending requests (`FetchPendingRequests`), the information is converted into a list of user maps and returned to the client as part of the `payload` in a message with `DONE` action.
- **UI:** The system provides feedback based on the operations (e.g., "Received a friend request from %s.").

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: Field definition mapping from `server/client.go`.

## Evidence index

- `server/client.go:27`: Field definition in `client` struct.
- `server/client.go:108`: Field modification (adding) in `SendFriendRequest`.
- `server/client.go:126`: Field modification (removing) in `AcceptFriendRequest`.
- `server/client.go:160`: Field modification (removing) in `DeleteFriendRequest`.
- `server/client.go:185`: Field modification (removing) in `RejectFriendRequest`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



