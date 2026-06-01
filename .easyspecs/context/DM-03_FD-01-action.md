# Field FD-01 — Action

**Entity:** DM-03 · **File:** DM-03_FD-01-action.md

## Summary
The `Action` field defines the type of action performed by the `Message` object. It determines how the payload of the message should be interpreted.

## Type and constraints
- **Type:** `ActionType` (an alias for `int`).
- **Constraints:** It is an enumeration defined using Go `iota`, representing specific action types such as `SIGN_UP`, `LOGIN`, `LOGOUT`, `CREATE_ROOM`, etc.

## Default and nullability
- Not explicitly defined in the struct. Assuming default `int` value (0, which corresponds to `SIGN_UP`).

## Validation
- Validation logic not directly inspected in `protocol/message.go`, but as a struct field, it is part of the `Message` serialization/deserialization process.

## Privacy / sensitivity
- Potentially sensitive as it indicates user intent (e.g., login, create room).

## Representation in API and UI
- Represented in JSON as `"action"`.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from protocol/message.go.

## Evidence index

- `protocol/message.go:5`: Definition of `ActionType`.
- `protocol/message.go:7-52`: Enum constants for `ActionType`.
- `protocol/message.go:55`: Usage of `ActionType` in `Message` struct.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Message](./DM-03-message.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



