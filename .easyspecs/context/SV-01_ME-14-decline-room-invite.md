# Method ME-14 — DeclineRoomInvite

**Service:** SV-01 · **File:** SV-01_ME-14-decline-room-invite.md

## Summary

Allows a user to decline an invitation to join a room.

## Operation

`DeclineRoomInvite(invitee *client, roomID, roomName string)`

This method first resolves the room using `roomID` and `roomName`. If found, it calls the `DeclineRoomInvite` method on the room instance, passing the `invitee` (client).

## Request / inputs

- `invitee`: The client declining the invite.
- `roomID`: The ID of the room.
- `roomName`: The name of the room.

## Response / outputs

- If successful, the invitee receives a success message confirming the decline.
- The room owner is notified that the invitation was declined.
- If the invite is not found or other errors occur, an error message is sent to the invitee.

## Auth and permissions

- Requires that the invitee has an active pending invitation to the specified room.

## Idempotency and concurrency

- The operation relies on `r.mu.Lock()` to ensure thread safety when modifying the room's `invites` map.

## Errors

- Errors are returned to the client using `invitee.err(err)`.
- Errors include:
  - Room not found (from `resolveRoom`).
  - Invite not found for the user in the room.
  - User is already a member of the room.

## Implementation notes

- The `DeclineRoomInvite` method is defined in `server/server.go` and delegates to `server/rooms.go`.

## Revision

- Initial draft: contract and handler from implementation.
- Refined evidence index and ensured implementation mapping.

## Evidence index

- `server/server.go:253-261`: Implementation of `server.DeclineRoomInvite`.
- `server/rooms.go:338-363`: Implementation of `room.DeclineRoomInvite`.
- `server/utils.go:211`: Implementation of `server.resolveRoom`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



