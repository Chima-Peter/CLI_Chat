# Field FD-08 — my_rooms

**Entity:** DM-02 · **File:** DM-02_FD-08-my-rooms.md

## Summary
The `my_rooms` field is a mapping of room IDs to `room` objects, representing the set of rooms that the client is currently a member of.

## Type and constraints
- **Type:** `map[string]*room`
- **Definition:** `server/client.go:25`
- **Constraints:** The key is a string (Room ID) and the value is a pointer to a `room` entity.

## Default and nullability
- **Default:** As a Go map, it defaults to `nil` if not initialized. It is typically initialized to an empty map when a client is created.

## Validation
- Access to this field is protected by a mutex `cl.mu` to ensure thread safety during concurrent read/write operations within the `client` structure.

## Privacy / sensitivity
- This field contains sensitive membership data for the client.

## Representation in API and UI
- Not directly serialized to JSON in the provided `client` struct methods, but used internally to track and manage the client's room memberships.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from schema/types.

## Evidence index

- `server/client.go:25`: Field definition `my_rooms map[string]*room`.
- `server/client.go:34`: Mutex definition `mu sync.RWMutex` used for synchronization of client fields.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



