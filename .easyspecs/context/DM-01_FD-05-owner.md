# Field FD-05 — owner

**Entity:** DM-01 · **File:** DM-01_FD-05-owner.md

## Summary
The `owner` field represents the `client` who created or holds administrative rights for the room.

## Type and constraints
- **Type:** `*client` (pointer to a `client` struct)
- **Source:** `server/rooms.go:15`

## Default and nullability
The field is a pointer, implying it can be `nil`. However, in practice, a room should have an owner upon creation to manage administrative tasks.

## Validation
The `owner` field is used for authorization in administrative actions:
- **EditRoom:** The system verifies `if r.owner.id != cl.id` to ensure only the owner can modify room settings (`server/rooms.go:158`).
- **DeleteRoom:** The system verifies `if r.owner.id != cl.id` to ensure only the owner can delete the room (`server/rooms.go:222`).

## Privacy / sensitivity
The `owner` field links a room to a specific `client` (user).

## Representation in API and UI
The field is used internally by the server logic to perform authorization checks.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/rooms.go.

## Evidence index

- `server/rooms.go:15`: Field definition as `*client`.
- `server/rooms.go:158`: Authorization check in `EditRoom`.
- `server/rooms.go:222`: Authorization check in `DeleteRoom`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Room](./DM-01-room.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



