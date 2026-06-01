# Field FD-06 — password

**Entity:** DM-01 · **File:** DM-01_FD-06-password.md

## Summary

The `password` field stores the optional password used to secure a private room. It is a string associated with an instance of `room`.

## Type and constraints

- **Type:** `string`
- **Constraints:** The password is stored after being trimmed of leading and trailing whitespace using `strings.TrimSpace()`. If the resulting password is empty, the room is considered public, and `is_private` is not set to true.

## Default and nullability

- **Default:** By default, a new room is public, and the password field is an empty string.

## Validation

- The password is validated during room creation or update via `SetRoomPassword`. An empty password (after trimming) indicates the room is public.
- When joining a private room, the provided password is compared against the stored password in `JoinRoomWithPassword`.

## Privacy / sensitivity

- The password is used to restrict access to the room. It is essential for protecting private rooms.

## Representation in API and UI

- The password is not directly exposed in the room information payload.
- It is requested from the user via the `GET_ROOM_PASSWORD` command when a user attempts to join a private room.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/rooms.go and fields list.

## Evidence index

- `server/rooms.go:16` (definition of password field in `room` struct)
- `server/rooms.go:62-74` (setting room password in `SetRoomPassword` method)
- `server/rooms.go:118-155` (validating password in `JoinRoomWithPassword` method)

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Room](./DM-01-room.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



