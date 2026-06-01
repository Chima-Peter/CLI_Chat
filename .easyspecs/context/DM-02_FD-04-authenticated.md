# Field FD-04 — authenticated

**Entity:** DM-02 · **File:** DM-02_FD-04-authenticated.md

## Summary

The `authenticated` field tracks the authentication status of a client session within the system. It determines whether a client has successfully logged in and is authorized to perform restricted actions.

## Type and constraints

- **Type:** `bool`
- **Constraints:** None (boolean flag)

## Default and nullability

- **Default:** `false` upon new connection initialization.
- **Nullability:** Not nullable (primitive boolean).

## Validation

The field itself is not validated against a schema, but its state is checked by the `dispatchMessage` handler in `server/handle_conn.go` to enforce restricted access. The `handleLogin` logic updates this flag upon successful nickname validation and uniqueness verification.

## Privacy / sensitivity

This field indicates the authentication status of a user session. It is not PII itself but is critical for session security.

## Representation in API and UI

The status is internal to the server's `client` structure and is used to control the flow of communication and action processing in the `dispatchMessage` loop.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/client.go and server/handle_conn.go.

## Evidence index

- `server/client.go:21`: Field definition in `client` struct.
- `server/handle_conn.go:21`: Initialized to `false` in `HandleConn`.
- `server/handle_conn.go:185`: Set to `true` upon successful login.
- `server/handle_conn.go:165`: Checked to prevent re-authentication of already authenticated clients.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



