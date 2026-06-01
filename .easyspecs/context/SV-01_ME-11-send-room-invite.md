# Method ME-11 — SendRoomInvite

**Service:** SV-01 · **File:** SV-01_ME-11-send-room-invite.md

## Summary

The `SendRoomInvite` method allows a room administrator to invite a new user to a room.

## Operation

`server.SendRoomInvite(cl *client, roomID, roomName, memberID, memberName string)`

This method resolves the target room and user, validates that the caller is the room administrator, ensures the user is not already a member, and then sends the invitation.

## Request / inputs

- `cl`: The client instance of the administrator initiating the invite.
- `roomID`: The unique identifier of the room.
- `roomName`: The name of the room (used for resolution if `roomID` is ambiguous or not provided).
- `memberID`: The unique identifier of the user to invite.
- `memberName`: The nickname of the user (used for resolution if `memberID` is ambiguous or not provided).

## Response / outputs

- Returns `void`. 
- Communicates errors back to the client (`cl`) via `cl.err(err)`.

## Auth and permissions

- Requires the caller (`cl`) to be the administrator of the room (`room_data.owner.id == cl.id`).

## Idempotency and concurrency

- Validates if the user is already a member of the room before sending the invite, preventing redundant invitations.

## Errors

- `err`: If room resolution fails.
- `err`: If the caller is not the room administrator.
- `err`: If user resolution fails.
- `err`: If the user is already a member of the room.

## Implementation notes

The method performs validation steps sequentially before calling `room_data.SendRoomInvite` to perform the actual invitation action.

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:201-226`: Implementation of the `SendRoomInvite` method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



