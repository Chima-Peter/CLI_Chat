# Service SV-01 — Server

**Slug:** server · **File:** SV-01-server.md

## Summary

The Server service acts as the central coordinator for the application, managing the state of chat rooms, connected clients, and mediating file transfer requests between peers. It maintains global state for the system and enforces room and user permissions.

## Responsibilities

- **Room Lifecycle Management:** Handles creation, editing, deletion of rooms, and managing room membership (joining/leaving).
- **User Management:** Manages friend requests, blocking users, and user status tracking.
- **Communication Mediation:** Orchestrates message passing within rooms and between friends.
- **File Transfer Facilitation:** Mediates file port requests and dispatches pending file transfers between clients.
- **Concurrency Control:** Ensures thread-safe access to shared state (rooms, clients, pending requests) using mutexes.

## Consumers

The primary consumers are connected clients, which interact with the server via an established connection and session (implemented in `client/client.go`).

## Public surface

The service exposes a set of methods for clients to interact with the system, primarily defined on the `server` struct in `server/server.go`. Key methods include:
- `CreateRoom`
- `JoinRoom` / `JoinRoomWithPassword`
- `LeaveRoom` / `DeleteRoom`
- `SendRoomMessage` / `SendContextMessage`
- `SendFriendRequest` / `AcceptFriendRequest` / `RejectFriendRequest`
- `requestFriendFilePort` / `handleFilePortListening`

## Dependencies

- **`sync`**: Utilized for mutexes (`sync.RWMutex`) to handle concurrent access to shared resources.
- **`github.com/google/uuid`**: Used for generating unique identifiers for rooms and sessions.
- **`encoding/json`**: For serializing and deserializing payloads.

## Error model

Errors are typically returned to the connected client directly via the client's own error handling mechanism (`cl.err(error)`). The server does not necessarily return errors to the caller but triggers callbacks on the client connection.

## Operational notes

The server holds the entire state of rooms and clients in memory using maps protected by a global `sync.RWMutex`. This implies that the server is currently designed to operate as a single-instance process.

## Revision

- Initial draft: Defined responsibilities, consumers, and summarized API surface from `server/server.go`.

## Evidence index

- `server/server.go:19-693`: Full `server` struct and methods implementation.
- `server/server.go:19-24`: `server` struct definition with shared state (rooms, clients, pending file requests, mutex).
- `server/server.go:36-74`: `CreateRoom` implementation.
- `server/server.go:403-468`: File transfer related methods (`requestFriendFilePort`, `handleFilePortListening`).

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [CreateRoom](./SV-01_ME-01-create-room.md)
- [SetRoomPassword](./SV-01_ME-02-set-room-password.md)
- [EditRoom](./SV-01_ME-03-edit-room.md)
- [JoinRoom](./SV-01_ME-04-join-room.md)
- [JoinRoomWithPassword](./SV-01_ME-06-join-room-with-password.md)
- [LeaveRoom](./SV-01_ME-07-leave-room.md)
- [DeleteRoom](./SV-01_ME-08-delete-room.md)
- [DeleteRoomMember](./SV-01_ME-09-delete-room-member.md)
- [GetRoomMembers](./SV-01_ME-10-get-room-members.md)
- [SendRoomInvite](./SV-01_ME-11-send-room-invite.md)
- [SeePendingRoomInvites](./SV-01_ME-12-see-pending-room-invites.md)
- [AcceptRoomInvite](./SV-01_ME-13-accept-room-invite.md)
- [DeclineRoomInvite](./SV-01_ME-14-decline-room-invite.md)
- [ListPublicRooms](./SV-01_ME-15-list-public-rooms.md)
- [SwitchContext](./SV-01_ME-16-switch-context.md)
- [SendContextMessage](./SV-01_ME-17-send-context-message.md)
- [requestFriendFilePort](./SV-01_ME-18-request-friend-file-port.md)
- [handleFilePortListening](./SV-01_ME-19-handle-file-port-listening.md)
- [dispatchPendingFileRequests](./SV-01_ME-20-dispatch-pending-file-requests.md)
- [sendFilePortToSender](./SV-01_ME-21-send-file-port-to-sender.md)
- [SendRoomMessage](./SV-01_ME-22-send-room-message.md)
- [ListMyRooms](./SV-01_ME-23-list-my-rooms.md)
- [ListMyRoomInvites](./SV-01_ME-24-list-my-room-invites.md)
- [SendFriendRequest](./SV-01_ME-25-send-friend-request.md)
- [AcceptFriendRequest](./SV-01_ME-26-accept-friend-request.md)
- [RejectFriendRequest](./SV-01_ME-27-reject-friend-request.md)
- [CancelFriendRequest](./SV-01_ME-28-cancel-friend-request.md)
- [SeePendingFriendRequests](./SV-01_ME-29-see-pending-friend-requests.md)
- [SeeSentFriendRequests](./SV-01_ME-30-see-sent-friend-requests.md)
- [GetFriends](./SV-01_ME-31-get-friends.md)
- [DeleteFriend](./SV-01_ME-32-delete-friend.md)
- [MessageFriend](./SV-01_ME-33-message-friend.md)
- [BlockUser](./SV-01_ME-34-block-user.md)
- [UnblockUser](./SV-01_ME-35-unblock-user.md)
- [GetUserStatus](./SV-01_ME-36-get-user-status.md)
- [LogUserOut](./SV-01_ME-37-log-user-out.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



