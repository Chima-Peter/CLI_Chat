# Field FD-03 — Payload

**Entity:** DM-03 · **File:** DM-03_FD-03-payload.md

## Summary

The `Payload` field holds the raw JSON payload of a message. It is part of the `Message` struct defined in the protocol, allowing for flexible message structures depending on the specific `Action` being performed.

## Type and constraints

The `Payload` field is defined as `json.RawMessage` in the Go `Message` struct (protocol/message.go:57). In Go, `json.RawMessage` is a `[]byte` used to delay JSON decoding or pre-encode a specific JSON structure. This means the field, when transmitted as JSON, must contain a valid JSON value (object, array, string, number, boolean, or null).

## Default and nullability

If the `Message` struct is initialized without explicitly setting `Payload`, the `json.RawMessage` field will be `nil`.

## Validation

The `Payload` field itself is not validated by the `Message` struct definitions beyond being required to be syntactically correct JSON if provided during unmarshaling. Semantic validation of the payload contents is dependent on the handler logic processing the `Action` specified in the message.

## Privacy / sensitivity

This field likely contains sensitive information or user data, such as message contents, file metadata, or authentication details, depending on the `Action` type. It should be treated as highly sensitive data according to the privacy requirements of the application.

## Representation in API and UI

In the JSON representation of the `Message` struct, this field is represented by the `payload` key.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from protocol/message.go.

## Evidence index

- `protocol/message.go:57`: Field definition of `Payload` as `json.RawMessage`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Message](./DM-03-message.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



