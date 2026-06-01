# Field FD-03 — nick

**Entity:** DM-02 · **File:** DM-02_FD-03-nick.md

## Summary

The `nick` field represents the unique nickname assigned to a client. It is used for identifying the user within the system, in room memberships, friend lists, and messaging.

## Type and constraints

- **Type:** `string` (Go)
- **Constraints:**
  - Must not be empty.
  - Must not be "anonymous" (case-insensitive).
  - Must be unique among all connected clients.

## Default and nullability

- **Default:** Initialized to an empty string `""` upon client connection.
- **Nullability:** Not nullable; it is a required `string` type.

## Validation

Validation occurs during the login process in the server:
- Nicknames are trimmed of surrounding whitespace using `strings.TrimSpace`.
- Empty strings are rejected.
- "anonymous" (case-insensitive) is a restricted nickname.
- Nickname uniqueness is enforced against currently connected clients using `s.isNickTaken` (which uses `strings.EqualFold` for case-insensitive comparison).

## Privacy / sensitivity

- The nickname is publicly visible to other users within the same room or as part of friend interactions.

## Representation in API and UI

- The nickname is used in UI prompts and log messages to identify users.
- It is sent in JSON payloads during communication (e.g., login, friend requests).

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/client.go and validation logic from server/handle_conn.go.

## Evidence index

- server/client.go:20 (Field definition in `client` struct)
- server/handle_conn.go:164-184 (Login handling and nickname validation logic)
- server/utils.go:116-124 (`isNickTaken` implementation)

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



