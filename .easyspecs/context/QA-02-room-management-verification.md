# Feature QA-02 — Room management verification

**Slug:** room-management-verification  
**Output file:** QA-02-room-management-verification.md

## Summary

This feature provides the capability to manage chat rooms, including room creation, membership handling, and configuration management. Users can join rooms, manage membership status, set room passwords, and administrators can modify room settings or delete rooms.

## Scope

In scope:
- Room creation and joining.
- Membership management (invites, joining, leaving).
- Room configuration (name, max size, privacy/passwords).
- Room deletion by owners.

Out of scope:
- Persistent storage of room configurations (in-memory only).

## Functional behaviour

The room management system allows users to interact with chat rooms:
- **Joining:** Users can join rooms. If a room is private, a password is required `server/rooms.go:118-155`.
- **Configuration:** Room owners can rename the room and set the maximum number of members `server/rooms.go:157-202`.
- **Membership:** Users can join, leave `server/rooms.go:204-214`, and be invited by the room owner `server/rooms.go:279-289`.
- **Deletion:** Room owners can delete the room, which removes all current members `server/rooms.go:216-248`.

## Technical design

The `room` structure in `server/rooms.go` acts as the primary data model for chat rooms.
- **Concurrency:** Thread safety is managed using `sync.RWMutex` to protect room data during concurrent operations `server/rooms.go:19`.
- **Membership:** `members` and `invites` are managed via maps, keyed by client ID `server/rooms.go:17-18`.

## Entry points

- `JoinRoom(cl *client)`: Handles initial joining logic `server/rooms.go:76-116`.
- `JoinRoomWithPassword(cl *client, password string)`: Handles joining private rooms `server/rooms.go:118-155`.
- `EditRoom(...)`: Modifies room settings `server/rooms.go:157-202`.
- `DeleteRoom(...)`: Handles room removal `server/rooms.go:216-248`.

## Dependencies

- `server` package for server-wide state and user management.
- `sync` for synchronization primitives.

## Open questions

None.

## Revision

- Initial draft: scope and behaviour from implementation reads.

## Evidence index

- `server/rooms.go:10-20`: Definition of the `room` struct.
- `server/rooms.go:76-116`: `JoinRoom` implementation.
- `server/rooms.go:157-202`: `EditRoom` implementation.
- `server/rooms.go:216-248`: `DeleteRoom` implementation.
- `server/rooms.go:279-363`: Invite-related methods implementation.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Join room](./QA-02_TC-01-join-room.md)
- [Join private room](./QA-02_TC-02-join-private-room.md)
- [Edit room](./QA-02_TC-03-edit-room.md)
- [Leave room](./QA-02_TC-04-leave-room.md)
- [Delete room](./QA-02_TC-05-delete-room.md)
- [Invite user to room](./QA-02_TC-06-invite-user-to-room.md)
- [View pending room invites](./QA-02_TC-07-view-pending-room-invites.md)
- [Accept room invite](./QA-02_TC-08-accept-room-invite.md)
- [Decline room invite](./QA-02_TC-09-decline-room-invite.md)
- [Set room password](./QA-02_TC-10-set-room-password.md)
- [Fetch room members](./QA-02_TC-11-fetch-room-members.md)
- [Broadcast message](./QA-02_TC-12-broadcast-message.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



