# Field FD-09 — friends

**Entity:** DM-02 · **File:** DM-02_FD-09-friends.md

## Summary

The `friends` field represents the set of users that a client is currently friends with. It is implemented as a map of strings to empty structs (`map[string]struct{}`), effectively acting as a set of unique friend identifiers.

## Type and constraints

- **Type:** `map[string]struct{}`
- **Constraints:**
    - Keys are `string` representing unique friend identifiers.
    - Values are `struct{}` (empty struct), used to represent the existence of a friend relationship in the set.
    - Thread safety: Access to this field should be protected by the client's mutex (`mu`), as inferred from its usage in `server/client.go`.

## Default and nullability

- **Default:** Initialized to an empty map `make(map[string]struct{})` when a new client is created in `server/handle_conn.go:27`.
- **Nullability:** The field itself is not a pointer, so it is never `nil` after initialization.

## Validation

- Relationships are validated during friend request handling and removals.
- The `hasID` helper function is commonly used to check for the existence of a specific ID in this map.

## Privacy / sensitivity

- This field contains information about the user's social connections, which may be considered private.

## Representation in API and UI

- The set of friends is returned in API responses, such as `GET_FRIENDS` protocol commands, often as a list of nicknames or IDs depending on the context of the requested operation.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/client.go and analysis of usage.

## Evidence index

- `server/client.go:26`: Field definition in `client` struct.
- `server/handle_conn.go:27`: Initialization of the `friends` map.
- `server/client.go:125,130`: Adding friends to the map.
- `server/client.go:248,252`: Removing friends from the map.
- `server/client.go:86,240,273`: Usage of `hasID` to check existence in the `friends` map.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



