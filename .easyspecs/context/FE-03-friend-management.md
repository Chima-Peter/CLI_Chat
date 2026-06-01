# Feature FE-03 — Friend management

**Slug:** friend-management
**Output file:** FE-03-friend-management.md

## Summary
The Friend management feature enables users to maintain a list of friends for direct messaging, manage incoming/outgoing friend requests, and block/unblock other users. It facilitates social interactions within the CLI chat system.

## Scope
- **In scope:**
  - Sending, accepting, rejecting, and canceling friend requests.
  - Listing existing friends and pending/sent requests.
  - Removing friends.
  - Messaging friends (direct messaging).
  - Blocking and unblocking users.
- **Out of scope:**
  - Persistence (all relationships are in-memory and volatile).

## Functional behaviour
Users interact with the friend management system via CLI commands, which trigger server handlers.
- Users can send a request to another user by nickname.
- Recipient users can accept, reject, or cancel requests.
- Friends can directly message each other after acceptance.
- Blocking a user prevents friend requests and direct messaging.

## Technical design
Friend relationships and request states are maintained in memory within the `client` struct (`server/client.go:17-35`).
- Relationships are stored in maps protected by `sync.RWMutex` (`server/client.go:34`).
- `friends`: `map[string]struct{}` storing IDs of friends (`server/client.go:26`).
- `pending_friend_requests`: `map[string]struct{}` for received requests (`server/client.go:27`).
- `sent_friend_request`: `map[string]struct{}` for sent requests (`server/client.go:28`).
- `server` methods (`server/server.go`) act as entry points for friend actions, validating user nicknames and triggering the corresponding `client` methods (`server/client.go`).

## Entry points
Entry points are defined in `server/server.go`:
- `SendFriendRequest`: `server/server.go:576-583`
- `AcceptFriendRequest`: `server/server.go:585-592`
- `RejectFriendRequest`: `server/server.go:594-601`
- `CancelFriendRequest`: `server/server.go:603-610`
- `SeePendingFriendRequests`: `server/server.go:612-617`
- `SeeSentFriendRequests`: `server/server.go:619-624`
- `GetFriends`: `server/server.go:626-631`
- `DeleteFriend`: `server/server.go:633-640`
- `MessageFriend`: `server/server.go:642-649`
- `BlockUser`: `server/server.go:651-658`
- `UnblockUser`: `server/server.go:660-667`

## Dependencies
- `server/server.go`: Core server infrastructure and request dispatching.
- `server/client.go`: Client state management and friend logic implementation.

## Open questions
None.

## Revision
- Initial draft: Feature description, scope, behaviour, design, and entry points documented based on `server/server.go` and `server/client.go` implementation.

## Evidence index
- `server/client.go:17-35`: Definition of `client` struct containing friend-related maps (`friends`, `pending_friend_requests`, `sent_friend_request`, `blocked_users`).
- `server/client.go:69-113`: Implementation of `SendFriendRequest` method in `client`.
- `server/client.go:115-136`: Implementation of `AcceptFriendRequest` method in `client`.
- `server/client.go:238-263`: Implementation of `DeleteFriend` method in `client`.
- `server/server.go:576-667`: Handler entry points for friend management actions in `server` struct.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Send friend request](./FE-03_UC-01.md)
- [Accept friend request](./FE-03_UC-02.md)
- [Reject friend request](./FE-03_UC-03.md)
- [Cancel sent friend request](./FE-03_UC-04.md)
- [List pending friend requests](./FE-03_UC-05.md)
- [List sent friend requests](./FE-03_UC-06.md)
- [List friends](./FE-03_UC-07.md)
- [Remove friend](./FE-03_UC-08.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



