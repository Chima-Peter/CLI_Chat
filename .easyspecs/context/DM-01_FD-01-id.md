# Field FD-01 — id

**Entity:** DM-01 · **File:** DM-01_FD-01-id.md

## Summary

The `id` field serves as the unique identifier for a room.

## Type and constraints

- **Type:** `string`
- **Constraints:** Must be unique for each room within the server.

## Default and nullability

- **Default:** None (must be assigned upon creation).
- **Nullability:** Non-nullable.

## Validation

Validation of room uniqueness based on `id` and `name` is handled in `server/server.go` (inferred).

## Privacy / sensitivity

This field is not sensitive, but is used in payload messages to identify the room context.

## Representation in API and UI

Used in JSON payloads for room-related operations (e.g., `Broadcast`, `JoinRoom`, `EditRoom`).

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/rooms.go.

## Evidence index

- `server/rooms.go:11`: `id` field defined as `string` within the `room` struct.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Room](./DM-01-room.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



