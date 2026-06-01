# Field FD-07 — members

**Entity:** DM-01 · **File:** DM-01_FD-07-members.md

## Summary

The `members` field tracks the set of clients currently joined to the room. It is implemented as a map, allowing efficient lookups and management of room participants.

## Type and constraints

- **Type:** `map[string]*client`
- **Key:** The client ID (string).
- **Value:** Pointer to the `client` structure representing the member.
- **Concurrency:** Access to this map is protected by the `sync.RWMutex` field (`mu`) within the `room` struct to ensure thread safety during concurrent operations.

## Default and nullability

The field is initialized as a map of clients upon the creation of a `room` struct. It is not nullable.

## Validation

- Room membership is enforced by `JoinRoom` and `JoinRoomWithPassword` methods, ensuring that a user cannot join a room they are already in.
- The `isFull()` method checks the length of the `members` map against `max_size` to prevent over-capacity.

## Representation in API and UI

The `members` field is primarily used internally within the `server` package for broadcasting messages (`Broadcast` method) and managing room-level operations like adding/removing members (`JoinRoom`, `DeleteMember`). While the raw map is not directly exposed in the API, the current room members' nicknames are accessible via the `FetchRoomMembers()` method.

## Revision

<!--
  Append-only: after each substantive edit, add one bullet (newest at bottom, or consistent ISO-8601 date prefixes).
-->

- Initial draft: field mapping from schema/types.

## Evidence index

- `server/rooms.go:17`: Definition of `members` as `map[string]*client` in `room` struct.
- `server/rooms.go:19`: Definition of `sync.RWMutex` for thread-safe access to `members`.
- `server/rooms.go:25`: Usage of `len(r.members)` in `isFull()` check.
- `server/rooms.go:111`: Addition of a client to `r.members` in `JoinRoom`.
- `server/rooms.go:255`: Removal of a client from `r.members` in `DeleteMember`.
- `server/rooms.go:272-273`: Iteration over `r.members` in `FetchRoomMembers()`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Room](./DM-01-room.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



