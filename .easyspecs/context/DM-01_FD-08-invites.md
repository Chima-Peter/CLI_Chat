# Field FD-08 — invites

**Entity:** DM-01 · **File:** DM-01_FD-08-invites.md

## Summary
The `invites` field stores a map of pending room invitations, where each entry maps a client's unique identifier (`string`) to a pointer to their `client` struct (`*client`).

## Type and constraints
- **Type:** `map[string]*client`
- **Constraints:** The key is a `string` (client ID) and the value is a pointer to a `client` structure. This field maintains references to pending invitations for room access.

## Default and nullability
- **Default:** If not explicitly initialized, the map may be `nil`.
- **Nullability:** In practice, this field is expected to hold a map instance.

## Validation
- Access to the map is protected by a mutex `r.mu` (RWMutex) to ensure thread-safe operations (adding, deleting, reading) as shown in `server/rooms.go:281-283`, `server/rooms.go:324-327`, and `server/rooms.go:355-357`.

## Privacy / sensitivity
- This field stores references to `client` objects who have been invited to the room. Access is managed through room methods and is restricted to the room's owner or via invite management actions.

## Representation in API and UI
- Not directly exposed to external APIs as a raw structure. Used internally to manage invite lifecycle (sending, accepting, declining).

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from schema/types.

## Evidence index

- `server/rooms.go:18`: Definition of the `invites` field in the `room` struct.
- `server/rooms.go:281-283`: Mutex-protected write operation adding an invite.
- `server/rooms.go:324-327`: Mutex-protected write operation removing an invite on acceptance.
- `server/rooms.go:355-357`: Mutex-protected write operation removing an invite on decline.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Room](./DM-01-room.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



