# Field FD-14 — online

**Entity:** DM-02 · **File:** DM-02_FD-14-online.md

## Summary
The `online` field represents the current connection status of a client. It is used to determine whether a client can receive messages or be interacted with by other users (e.g., friends).

## Type and constraints
- **Type:** `bool`
- **Definition:** Defined in `server/client.go` at line 31 as part of the `client` struct.
- **Constraints:** Access to this field is protected by the `client` struct's `sync.RWMutex` (`mu`) to ensure thread safety during concurrent read/write operations.

## Default and nullability
- **Default:** As a `bool` in Go, the zero value is `false`, implying that a newly initialized client is offline by default until explicitly set otherwise.
- **Nullability:** This is a primitive `bool` type in Go, so it is not nullable.

## Validation
- The field is updated via the `setOnline` method (lines 360-364), which locks the mutex before modifying the value.
- When reading the field (e.g., in `MessageFriend` or `GetUserStatus`), the `RWMutex.RLock()` is acquired to ensure a consistent read.

## Privacy / sensitivity
- This field indicates whether a user is currently connected to the server. This information is used for presence features, such as showing status in `GetUserStatus` and checking availability before messaging.

## Representation in API and UI
- The field is used in `GetUserStatus` to inform other users of the client's online status (sent via the `online` key in the JSON payload, line 389).
- It is checked in `MessageFriend` to prevent sending messages to offline users (line 291).

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/client.go.

## Evidence index

- `server/client.go:31`: Definition of `online` field in `client` struct.
- `server/client.go:360-364`: `setOnline` method used to update the field.
- `server/client.go:291-292`: Reading `online` status in `MessageFriend`.
- `server/client.go:372`: Reading `online` status in `GetUserStatus`.
- `server/client.go:389`: Including `online` in the payload for `GetUserStatus` API response.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



