# Field FD-01 — id

**Entity:** DM-02 · **File:** DM-02_FD-01-id.md

## Summary
The `id` field acts as a unique identifier for a `client` entity.

## Type and constraints
- **Type:** `string`
- **Constraints:** Must be unique for each client.

## Default and nullability
- **Default:** Unspecified in the struct, but Go string default is `""` (empty string).
- **Nullability:** Non-nullable (it's a `string`, not a pointer to string).

## Validation
- The `id` is used extensively in comparisons, particularly for friend requests (`server/client.go:70`), blocking users (`server/client.go:324`), and checking friendship/block status.

## Privacy / sensitivity
- The ID itself is likely used as an internal identifier; it might be exposed in API responses (e.g., `server/client.go:313`).

## Representation in API and UI
- Appears in internal message payloads (`server/client.go:313`).

## Revision

- Initial draft: field mapping from server/client.go.

## Evidence index

- `server/client.go:18` — Definition of `id` field in `client` struct.
- `server/client.go:70` — Use of `id` in `SendFriendRequest`.
- `server/client.go:313` — Use of `id` in `MessageFriend` payload.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



