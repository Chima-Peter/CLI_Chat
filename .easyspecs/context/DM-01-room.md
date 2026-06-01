# Entity DM-01 — Room

**Slug:** room · **File:** DM-01-room.md

## Summary

The `Room` entity acts as a central hub for users (represented by the `Client` entity) to exchange messages. It manages its own state, including membership, access control (private rooms/passwords), and configuration.

## Purpose and lifecycle

A `Room` is created and managed by an owner (`Client`). Users can join rooms, optionally requiring a password if the room is private. Rooms can be edited by the owner to change the name or capacity. A room can be deleted by its owner, which removes all members and cleans up state in both the room and the associated clients.

## Invariants

- A room has a unique `id` and `name`.
- A room may have a capacity limit (`max_size`).
- If `is_private` is `true`, users must provide the correct `password` to join.
- Only the `owner` of the room can edit or delete it.
- A user cannot be a member of a full room.
- Membership and state changes are protected by a `sync.RWMutex` to ensure thread safety.

## Storage mapping

The `Room` entity is an in-memory structure defined in the `server` package.

```go
type room struct {
	id         string
	name       string
	max_size   int
	is_private bool
	owner      *client
	password   string
	members    map[string]*client
	invites    map[string]*client
	mu         sync.RWMutex
}
```

## Fields overview

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | `string` | Unique identifier for the room. |
| `name` | `string` | Display name of the room. |
| `max_size` | `int` | Maximum allowed members. |
| `is_private` | `bool` | Flag for password protection. |
| `owner` | `*client` | The client who created/owns the room. |
| `password` | `string` | Password for private rooms. |
| `members` | `map[string]*client` | Map of active members. |
| `invites` | `map[string]*client` | Map of pending invites. |
| `mu` | `sync.RWMutex` | Mutex for thread-safe access. |

## Relationships overview

- **Owner**: A `Room` has one `owner` of type `*client`.
- **Members**: A `Room` has many `members` of type `*client`.
- **Invites**: A `Room` can have many pending `invites` to `*client`s.

## Revision

- Initial draft: defined entity lifecycle, invariants, and storage mapping based on `server/rooms.go`.

## Evidence index

- `server/rooms.go:10-20`: `room` struct definition.
- `server/rooms.go:22-26`: `isFull()` invariant check.
- `server/rooms.go:76-116`: `JoinRoom()` logic handling privacy and capacity.
- `server/rooms.go:157-202`: `EditRoom()` logic for owner-only access and capacity validation.
- `server/rooms.go:216-248`: `DeleteRoom()` logic for owner-only access and cleanup.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [id](./DM-01_FD-01-id.md)
- [name](./DM-01_FD-02-name.md)
- [max_size](./DM-01_FD-03-max-size.md)
- [is_private](./DM-01_FD-04-is-private.md)
- [owner](./DM-01_FD-05-owner.md)
- [password](./DM-01_FD-06-password.md)
- [members](./DM-01_FD-07-members.md)
- [invites](./DM-01_FD-08-invites.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



