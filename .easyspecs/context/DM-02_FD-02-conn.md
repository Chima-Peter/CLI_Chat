# Field FD-02 — conn

**Entity:** DM-02 · **File:** DM-02_FD-02-conn.md

## Summary

The `conn` field represents the network connection for the client instance. It is used by the server to send and receive data from the connected client.

## Type and constraints

The `conn` field is of type `net.Conn` from the Go `net` package.

- **Type:** `net.Conn`
- **Location:** `server/client.go:19`

## Default and nullability

As it is a `net.Conn` (which is an interface type in Go), the default value is `nil` if the client is not connected.

## Validation

The `conn` object itself is not subject to structural validation via tags or explicit checks within the `client` struct definition. Its validity is managed by the Go `net` package when establishing and managing the TCP/network connection.

## Privacy / sensitivity

This field represents an active network connection object. While the connection itself is used to transmit data, it is an internal server-side representation of the socket.

## Representation in API and UI

This field is purely internal to the server's management of clients. It is not exposed directly in public-facing APIs or the user interface.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/client.go.

## Evidence index

- `server/client.go:19` (Definition of `conn` in `client` struct)

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



