# Entity DM-02 — Client

**Slug:** client · **File:** DM-02-client.md

## Summary

The `Client` entity represents a connected user within the chat system, tracking their connection state, authentication, presence, relationships (friends, blocks), and activity context (current room or friend conversation).

## Purpose and lifecycle

The `Client` is instantiated when a user connects to the server. It persists as long as the connection remains active, managing the user's state through a `sync.RWMutex` to ensure thread-safe access during concurrent operations like messaging, friend management, and status updates.

## Invariants

- Every `Client` has a unique ID and an associated `net.Conn`.
- The `mu` `sync.RWMutex` must be used to protect access to internal state fields (e.g., `friends`, `blocked_users`, `online` status).
- Friends relationships are bi-directional and managed in tandem between two `Client` instances.
- A user cannot block themselves.
- Friend requests, acceptance, and rejection operations involve synchronized updates across both the sender's and receiver's client states.

## Storage mapping

The `Client` is an in-memory entity defined in `server/client.go`. There is no direct persistent storage (e.g., database) mapping visible; state is transient and exists for the duration of the user's session.

## Fields overview

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | `string` | Unique identifier. |
| `conn` | `net.Conn` | Network connection. |
| `nick` | `string` | User nickname. |
| `authenticated` | `bool` | Authentication state. |
| `room` | `*room` | Current chat room, if any. |
| `current_context` | `string` | Current interaction context (`room`, `friend`, `none`). |
| `current_friend` | `*client` | Active friend conversation target, if any. |
| `my_rooms` | `map[string]*room` | Rooms joined by the user. |
| `friends` | `map[string]struct{}` | Set of friend user IDs. |
| `pending_friend_requests` | `map[string]struct{}` | Received friend requests. |
| `sent_friend_request` | `map[string]struct{}` | Sent friend requests. |
| `blocked_users` | `map[string]struct{}` | Users blocked by this client. |
| `room_invites` | `map[string]*room` | Received room invitations. |
| `online` | `bool` | Online/offline status. |
| `fileListenHost` | `string` | Host for file transfers. |
| `fileListenPort` | `int` | Port for file transfers. |
| `mu` | `sync.RWMutex` | Thread safety mechanism. |

## Relationships overview

- **Rooms**: A `Client` can be in a `room` (`DM-01`) and can have multiple `my_rooms`.
- **Friends**: `Client` instances can have friend relationships with other `Client` instances, managed through friend requests and blocks.

## Revision

- Initial draft: Entity lifecycle from `server/client.go`.

## Evidence index

- `server/client.go:17-35`: `client` struct definition.
- `server/client.go:58-60`: Use of `RLock` for thread-safe reading of `blocked_users`.
- `server/client.go:103-109`: Use of `Lock` for thread-safe update of friend request state.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [id](./DM-02_FD-01-id.md)
- [conn](./DM-02_FD-02-conn.md)
- [nick](./DM-02_FD-03-nick.md)
- [authenticated](./DM-02_FD-04-authenticated.md)
- [room](./DM-02_FD-05-room.md)
- [current_context](./DM-02_FD-06-current-context.md)
- [current_friend](./DM-02_FD-07-current-friend.md)
- [my_rooms](./DM-02_FD-08-my-rooms.md)
- [friends](./DM-02_FD-09-friends.md)
- [pending_friend_requests](./DM-02_FD-10-pending-friend-requests.md)
- [sent_friend_request](./DM-02_FD-11-sent-friend-request.md)
- [blocked_users](./DM-02_FD-12-blocked-users.md)
- [room_invites](./DM-02_FD-13-room-invites.md)
- [online](./DM-02_FD-14-online.md)
- [fileListenHost](./DM-02_FD-15-file-listen-host.md)
- [fileListenPort](./DM-02_FD-16-file-listen-port.md)
- [mu](./DM-02_FD-17-mu.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



