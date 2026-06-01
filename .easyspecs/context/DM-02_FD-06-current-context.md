# Field FD-06 — current_context

**Entity:** DM-02 · **File:** DM-02_FD-06-current-context.md

## Summary

The `current_context` field stores the current interaction context of a `client` within the system. It represents the state or focus of the client's current activity.

## Type and constraints

- **Type:** `string` (Go primitive type).

## Default and nullability

- **Default:** In Go, struct fields that are strings and not explicitly initialized default to an empty string (`""`).
- **Nullability:** As a `string` type in this Go struct, it is not nullable.

## Validation

No explicit validation logic is currently visible for this field within the `client` struct definition.

## Privacy / sensitivity

This field likely contains sensitive information related to user activity and current interaction state. Depending on the content stored, it could be considered PII or otherwise sensitive.

## Representation in API and UI

This field is part of the `client` structure in `server/client.go`. Its exposure via API endpoints or UI components is not currently documented in this view.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from `server/client.go`.

## Evidence index

- `server/client.go:23` (Field definition)

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



