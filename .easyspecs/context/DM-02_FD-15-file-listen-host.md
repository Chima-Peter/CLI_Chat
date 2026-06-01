# Field FD-15 — fileListenHost

**Entity:** DM-02 · **File:** DM-02_FD-15-file-listen-host.md

## Summary
The `fileListenHost` field stores the hostname or IP address on which the client is listening for incoming file transfers.

## Type and constraints
- **Type:** `string`
- **Constraints:** None explicitly defined in the struct, but expected to be a valid network hostname or IP address.

## Default and nullability
- **Default:** Go zero-value for `string` (empty string `""`).
- **Nullability:** Not a pointer, so cannot be `nil`, but can be an empty string.

## Validation
There is no explicit validation logic implemented for this field in `server/client.go` at the point of definition.

## Privacy / sensitivity
This field exposes network information about where the client is accepting file connections, which may be sensitive depending on the infrastructure (e.g., public vs. local IP).

## Representation in API and UI
The field is retrieved via `GetUserFilePort()` in `server/client.go`.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/client.go.

## Evidence index

- `server/client.go:32`: Field definition.
- `server/client.go:398`: Field usage in `GetUserFilePort`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



