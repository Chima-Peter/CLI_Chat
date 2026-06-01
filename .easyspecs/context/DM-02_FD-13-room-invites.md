# Field FD-13 — room_invites

**Entity:** DM-02 · **File:** DM-02_FD-13-room-invites.md

## Summary

The `room_invites` field maintains a mapping of room IDs to `room` objects, representing the pending room invitations received by a client.

## Type and constraints

The field is defined as a Go map: `map[string]*room`.

- Key: `string` (representing the room identifier)
- Value: `*room` (a pointer to a `room` struct)

## Default and nullability

As it is a map, the zero value is `nil` (or an empty map if initialized). In the context of the `client` struct, it is not explicitly initialized in the struct definition itself, suggesting it should be initialized before use to avoid nil map panics.

## Validation

No specific validation logic beyond map operations (`make`, `delete`, insertion).

## Privacy / sensitivity

This field tracks pending invitations to rooms, which could be considered sensitive as it reveals information about potential room access.

## Representation in API and UI

It likely represents pending invites that the user can accept or decline.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/client.go.

## Evidence index

- `server/client.go:30`: Field definition `room_invites map[string]*room`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



