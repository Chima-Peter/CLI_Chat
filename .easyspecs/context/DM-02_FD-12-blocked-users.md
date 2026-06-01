# Field FD-12 — blocked_users

**Entity:** DM-02 · **File:** DM-02_FD-12-blocked-users.md

## Summary
The `blocked_users` field maintains a set of user identifiers representing clients that have been blocked by this client. This allows for blocking communication or interactions with specific users.

## Type and constraints
Type: `map[string]struct{}`.
The map keys are strings representing user identifiers. The values are empty structs (`struct{}`), which is a common Go idiom for implementing a set in a map.

## Default and nullability
Initialized as an empty map upon client creation. It is not nullable.

## Validation
N/A (Managed by application logic).

## Privacy / sensitivity
PII/Sensitive: Contains identifiers of blocked users. Access should be restricted to the owning client's context.

## Representation in API and UI
Internal to the client session. Not typically exposed directly to other clients.

## Revision

- Initial draft: field mapping from schema/types.

## Evidence index

- Definition: `server/client.go:29`
- Initialization: `server/handle_conn.go:30`
- Usage for blocking check: `server/client.go:59`, `server/client.go:63`, `server/client.go:330`
- Blocking logic: `server/client.go:338`
- Unblocking logic: `server/client.go:354`

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



