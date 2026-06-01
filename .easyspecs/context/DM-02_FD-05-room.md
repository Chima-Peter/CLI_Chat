# Field FD-05 — room

**Entity:** DM-02 · **File:** DM-02_FD-05-room.md

## Summary
The `room` field represents the room that the client is currently in. It is a pointer to a `room` struct, allowing the client to be associated with a specific room instance.

## Type and constraints
- **Type:** `*room`
- **Definition:** The `room` struct is defined in `server/rooms.go:10`.

## Default and nullability
- **Nullable:** Yes, this is a pointer field, and it will be `nil` if the client is not currently in a room.

## Validation
- There is no specific validation logic shown in `server/client.go` at the definition site. It is managed by room-related services when a client joins or leaves a room.

## Privacy / sensitivity
- Not explicitly PII, but it reveals the client's current location within the system, which is part of their active session state.

## Representation in API and UI
- This field represents the client's current room state in the backend. Its serialization would depend on how the `client` struct is converted for API responses.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from schema/types.

## Evidence index

- `server/client.go:22`: Definition of `room` field in `client` struct.
- `server/rooms.go:10`: Definition of `room` struct type.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



