# Field FD-02 — name

**Entity:** DM-01 · **File:** DM-01_FD-02-name.md

## Summary

The `name` field represents the unique identifier (display name) for a room. It is used for identifying rooms, displaying room information in messages, and ensuring room name uniqueness within the server.

## Type and constraints

- **Type:** `string`
- **Constraint:** Must be unique across all rooms within the server instance.

## Default and nullability

- **Default:** Not explicitly defined at the struct level in Go, but initialized upon room creation.
- **Nullability:** Not nullable; it is a required `string` field.

## Validation

Room name validation is enforced during room editing:
- The name is trimmed of leading/trailing whitespace (`strings.TrimSpace`).
- The system checks if the proposed name is already taken by another room using `s.isRoomNameTaken(newName, r.id)`.

## Privacy / sensitivity

This field is not considered sensitive. It is public information used to identify and join rooms.

## Representation in API and UI

The room name is included in various API/message payloads, for example:
- When joining a room: `"room": r.name`.
- In room broadcast messages: `"room": r.name`.
- In room settings updates: `"room": r.name`.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from schema/types.

## Evidence index

- `server/rooms.go:12`: Definition of `name` field as `string` in `room` struct.
- `server/rooms.go:183-189`: Validation logic checking for room name uniqueness (`isRoomNameTaken`) before updating `r.name`.
- `server/rooms.go:163`: Trimming of whitespace in `EditRoom`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Room](./DM-01-room.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



