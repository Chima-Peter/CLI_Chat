# Field FD-16 — fileListenPort

**Entity:** DM-02 · **File:** DM-02_FD-16-file-listen-port.md

## Summary
The `fileListenPort` field represents the port on which the client is listening for incoming file transfers.

## Type and constraints
- **Type:** `int`
- **Constraints:** Must be a valid network port number (though not explicitly constrained in the struct definition, it is intended to represent a network port).

## Default and nullability
- **Default:** The zero value for `int`, which is `0`.
- **Nullability:** Not nullable as it is a primitive `int`.

## Validation
There is no explicit validation logic for the port number within the struct definition.

## Privacy / sensitivity
This field is part of the client's connectivity information for file transfers.

## Representation in API and UI
The port is retrieved via the `GetUserFilePort` method, which also returns the `fileListenHost`.

```go
func (cl *client) GetUserFilePort() (string, int) {
	cl.mu.RLock()
	fileListenHost := cl.fileListenHost
	fileListenPort := cl.fileListenPort
	cl.mu.RUnlock()

	return fileListenHost, fileListenPort
}
```

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from server/client.go.

## Evidence index

- `server/client.go:33` - Definition of `fileListenPort` field.
- `server/client.go:396-403` - `GetUserFilePort` method using the field.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



