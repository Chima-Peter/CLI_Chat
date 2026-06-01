# Method ME-13 — AcceptRoomInvite

**Service:** SV-01 · **File:** SV-01_ME-13-accept-room-invite.md

## Summary

The `AcceptRoomInvite` method allows a user to accept a pending invitation to a specific chat room. It validates the invite status, checks if the user is already a member, verifies room capacity, and completes the joining process by updating the client and room states, notifying the user, and broadcasting the join event to other members.

## Operation

The operation is handled in two stages:
1.  **Handler**: `server/server.go:243-251` resolves the room and invokes the room's business logic.
2.  **Logic**: `server/rooms.go:302-336` performs the validation and state transitions.

## Request / inputs

- `invitee`: `*client` - The client attempting to accept the invite.
- `roomID`: `string` - The ID of the room to join.
- `roomName`: `string` - The name of the room to join (used for resolution).

## Response / outputs

- None (void).
- **Side effects**:
    - Updates `invitee`'s room memberships and clears the invitation.
    - Updates `room`'s membership list.
    - Sends success message to the `invitee`.
    - Broadcasts "just joined" message to the room members.
    - Sends error messages to `invitee` on failure.

## Auth and permissions

- Requires that a pending invite for the user exists for the given `roomID`.

## Idempotency and concurrency

- The `room.AcceptRoomInvite` method uses a mutex `r.mu` to protect modifications to the room's `invites` and `members` collections (`server/rooms.go:324-327`).

## Errors

- "You have not been invited to this room" (invite not found).
- "You are already a part of this room" (already a member).
- "Room is full (max %d members)." (capacity reached).

## Implementation notes

- The handler delegates to the `room` struct's `AcceptRoomInvite` method.
- The state is updated both on the `invitee` client object and the `room` object.

## Revision

- Initial draft: contract, handler, and business logic from implementation.

## Evidence index

- `server/server.go:243-251`: Implementation of the `AcceptRoomInvite` handler.
- `server/rooms.go:302-336`: Implementation of the `AcceptRoomInvite` business logic.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



