# Feature FE-02 — Room management

**Slug:** room-management  
**Output file:** FE-02-room-management.md

## Summary

The Room Management feature provides core functionality for users to create, interact with, and govern chat rooms. It encompasses lifecycle management (creation, deletion, configuration), membership management (joining, leaving, inviting), and communication within rooms (broadcasting).

## Scope

- **In scope:**
    - Room creation and configuration (name, max size, privacy/passwords).
    - Room joining (public, private with password).
    - Room administration by the owner (editing settings, deleting).
    - Room membership management (invitations, accepting/declining invites, removing members).
    - Message broadcasting to room members.
- **Out of scope:**
    - Global room listing (handled by server infrastructure).
    - Client authentication (FE-01).

## Functional behaviour

The feature allows users to:
- **Join/Leave:** Join rooms (with password validation if private) and leave them.
- **Broadcast:** Send messages to all members of a room.
- **Manage Members:** Invite other users, accept/decline invites, and manage membership.
- **Administer:** Owners can rename rooms, change maximum size, set passwords, and delete rooms.

## Technical design

Room functionality is centralized in the `room` struct, which uses `sync.RWMutex` to manage thread-safe access to membership maps and room settings.

- **Storage:** Rooms maintain a `members` map and an `invites` map of `*client` objects.
- **Concurrency:** All sensitive operations on room state (membership, settings) are protected by a `sync.RWMutex`.
- **Communication:** The `Broadcast` method handles delivering messages to all members, providing appropriate context if the sender is or is not in the room.

## Entry points

Room management operations are invoked via the server's command handlers, which trigger the corresponding methods on the `room` object (e.g., `JoinRoom`, `EditRoom`, `DeleteRoom`).

## Dependencies

- **`client` struct:** Represents connected users participating in rooms.
- **`server` struct:** Manages the collection of all rooms and facilitates cross-room or global administrative operations (like room deletion).
- **`sync`, `maps`, `fmt`, `strings`:** Standard Go library for thread-safe state, data manipulation, and string handling.

## Open questions

None.

## Revision

- Initial draft: scope and behaviour from implementation reads in `server/rooms.go`.

## Evidence index

- `server/rooms.go:10-20`: Definition of `room` struct with thread-safe mechanisms and membership tracking.
- `server/rooms.go:28-60`: `Broadcast` method logic for delivering room messages.
- `server/rooms.go:76-116`: `JoinRoom` implementation handling public and private room entry.
- `server/rooms.go:157-202`: `EditRoom` implementation allowing owners to manage room settings.
- `server/rooms.go:216-248`: `DeleteRoom` logic for administrative room removal.
- `server/rooms.go:279-300`: Invitation handling logic.
- `server/rooms.go:302-336`: `AcceptRoomInvite` handling room membership transitions.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Join room](./FE-02_UC-01.md)
- [Join private room](./FE-02_UC-02.md)
- [Edit room](./FE-02_UC-03.md)
- [Leave room](./FE-02_UC-04.md)
- [Delete room](./FE-02_UC-05.md)
- [Invite user to room](./FE-02_UC-06.md)
- [View pending room invites](./FE-02_UC-07.md)
- [Accept room invite](./FE-02_UC-08.md)
- [Decline room invite](./FE-02_UC-09.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



