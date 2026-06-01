# Field FD-03 — max_size

**Entity:** DM-01 · **File:** DM-01_FD-03-max-size.md

## Summary

The `max_size` field defines the maximum number of members allowed to join a room.

## Type and constraints

The field is of type `int` (`server/rooms.go:13`).

## Default and nullability

While not explicitly initialized in the struct definition beyond the default `0` value for `int` in Go, it is interpreted as:
- A `max_size` of `0` indicates no limit on the number of members.

## Validation

- **Negative value:** In `EditRoom`, `max_size` cannot be set to a negative value (`server/rooms.go:167-171`).
- **Current members:** If `max_size` is set to a value greater than `0`, it cannot be lower than the current number of members in the room (`server/rooms.go:172-176`).
- **Fullness:** A room is considered full when `max_size > 0` and the number of current members is greater than or equal to `max_size` (`server/rooms.go:25`).

## Privacy / sensitivity

Not sensitive.

## Representation in API and UI

Included in the room payload during room updates (`server/rooms.go:197`).

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/rooms.go.

## Evidence index

- `server/rooms.go:13`: Definition of `max_size` field.
- `server/rooms.go:25`: Usage in `isFull` logic.
- `server/rooms.go:167-176`: Validation logic in `EditRoom`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Room](./DM-01-room.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



