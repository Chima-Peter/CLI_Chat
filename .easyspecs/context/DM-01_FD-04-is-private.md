# Field FD-04 — is_private

**Entity:** DM-01 · **File:** DM-01_FD-04-is-private.md

## Summary
The `is_private` field indicates whether a room requires a password for users to join.

## Type and constraints
- **Type:** `bool` (Go `bool`)
- **Constraints:** This field is managed internally as part of the `room` struct. It is set to `true` when a password is assigned to the room.

## Default and nullability
- **Default:** Not explicitly initialized, defaults to the Go `bool` zero value (`false`), indicating a public room by default.
- **Nullability:** Non-nullable (`bool`).

## Validation
- When `SetRoomPassword` is called, if the provided password (after trimming) is not empty, `is_private` is set to `true` (server/rooms.go:70).
- If the password is empty (after trimming), the room effectively remains public.

## Privacy / sensitivity
- This field is not directly sensitive but controls access to the room by requiring a password, which is a security mechanism.

## Representation in API and UI
- The field is used in `JoinRoom` to determine if the user needs to provide a password before joining (server/rooms.go:99).

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/rooms.go and fields-list.

## Evidence index

- `server/rooms.go:14`: Field definition in `room` struct.
- `server/rooms.go:70`: Field set to `true` in `SetRoomPassword`.
- `server/rooms.go:99`: Field checked in `JoinRoom` to gate access.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Room](./DM-01-room.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



