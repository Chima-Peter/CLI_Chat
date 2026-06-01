# Field FD-17 — mu

**Entity:** DM-02 · **File:** DM-02_FD-17-mu.md

## Summary

The `mu` field provides a mutex for thread-safe synchronization of the `client` struct instance. It ensures that concurrent access to the client's state (such as friend lists, block lists, and connection status) is handled safely.

## Type and constraints

- **Type:** `sync.RWMutex`
- **Constraints:** None specific to the field itself, but it must be used to protect all access to other fields within the `client` struct that can be accessed concurrently.

## Default and nullability

- **Default:** A `sync.RWMutex` initialized to its zero value, which is valid for use as a read-write lock.
- **Nullability:** Not nullable; it is a struct instance.

## Validation

- **Validation:** No functional validation (the type ensures its own integrity). Its correct usage is validated implicitly by concurrent access patterns throughout the `server/client.go` file.

## Privacy / sensitivity

- **Sensitivity:** The mutex itself is not sensitive, but it protects sensitive data like friend lists, block lists, and room information within the `client` struct.

## Representation in API and UI

- **Representation:** Not exposed in the API or UI; it is an internal implementation detail for thread safety.

## Revision

- <!-- Initial draft: field mapping from schema/types. -->

## Evidence index

- `server/client.go:34`: Definition of the `mu` field as `sync.RWMutex` within the `client` struct.
- `server/client.go:58-64`: Usage of `cl.mu.RLock()` and `cl.mu.RUnlock()` to safely read the `blocked_users` map.
- `server/client.go:103-105`: Usage of `cl.mu.Lock()` and `cl.mu.Unlock()` to safely write to the `sent_friend_request` map.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Client](./DM-02-client.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



