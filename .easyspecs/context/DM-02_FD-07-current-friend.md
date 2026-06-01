# Field FD-07 — current_friend

**Entity:** DM-02 · **File:** DM-02_FD-07-current-friend.md

## Summary

The `current_friend` field is a pointer to a `client` struct, representing the specific peer with whom the client is currently in a direct, active friend-to-friend interaction context.

## Type and constraints

- **Type:** `*client` (pointer to a `client` struct, defined in `server/client.go`).
- **Constraint:** This field is part of the `client` struct, which is used for managing client-side state in the server.

## Default and nullability

- **Default:** The field is initialized to `nil` when a new client is created (see `server/handle_conn.go`).
- **Nullability:** It is nullable (and frequently checked for `nil` to determine if a direct friend context is active).

## Validation

The field's value is checked for `nil` before accessing its attributes (e.g., `friend.current_friend != nil`), ensuring safe interaction with the associated friend's state.

## Privacy / sensitivity

This field holds a reference to another client's internal representation, which may be sensitive information as it indicates active P2P interaction context.

## Representation in API and UI

It influences message routing and context-specific functionality, such as direct friend-to-friend file transfers or messaging.

## Revision

- Initial draft: field mapping from schema/types.

## Evidence index

- `server/client.go:24`: Definition of `current_friend` field as `*client` pointer.
- `server/handle_conn.go:24`: Initialization of `current_friend` to `nil` in the client structure.
- `server/server.go:330`: Resetting `current_friend` to `nil`.
- `server/server.go:355`: Setting `current_friend` to a specific `friend` client object.
- `server/server.go:379`: Retrieving the `current_friend` reference.
- `server/client.go:299-305`: Usage of `current_friend` for active interaction checks.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



