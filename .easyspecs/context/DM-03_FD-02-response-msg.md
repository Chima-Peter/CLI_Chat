# Field FD-02 — ResponseMsg

**Entity:** DM-03 · **File:** DM-03_FD-02-response-msg.md

## Summary
The `ResponseMsg` field is a `string` property within the `Message` struct, used to convey informational, success, or error messages from the server to the client.

## Type and constraints
- **Type:** `string`
- **JSON tag:** `response_msg`

## Default and nullability
- In the Go `Message` struct, it is a `string` and defaults to an empty string `""` when a `Message` is initialized.

## Validation
- No explicit structural validation is applied to this field in the `protocol` package; it is handled as a standard string.

## Privacy / sensitivity
- This field may contain user-visible information, status messages, or potential error descriptions which might be sensitive depending on the content (e.g., system paths or internal identifiers in error messages).

## Representation in API and UI
- The `ResponseMsg` is used in the communication protocol to pass feedback to the client. In the client implementation (`client/session.go`), it is often directly written to the display if present (`s.writeDisplayLocked(msg.ResponseMsg)`).

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from protocol/message.go.

## Evidence index

- `protocol/message.go:56`: Definition of `ResponseMsg` field in `Message` struct.
- `client/session.go:105-110`: Example usage where `ResponseMsg` is written to the display.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Message](./DM-03-message.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



