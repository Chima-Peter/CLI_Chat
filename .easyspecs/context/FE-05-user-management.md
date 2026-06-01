# Feature FE-05 — User management

**Slug:** user-management  
**Output file:** FE-05-user-management.md

## Summary

The User management feature enables users to control their interactions with others by blocking unwanted users, unblocking previously blocked users, and checking the block status of other users.

## Scope

- **In-scope:**
  - Blocking a user by nickname.
  - Unblocking a user by nickname.
  - Checking if a user is currently blocked.
- **Out-of-scope:**
  - Administrative blocking or banning of users.
  - Managing user authentication (FE-01).
  - Managing friend relationships (FE-03).

## Functional behaviour

Users interact with this feature via specific terminal commands:
- `/user/block <user>`: Adds a user to the current user's blocked list, preventing further interactions from them.
- `/user/unblock <user>`: Removes a user from the current user's blocked list, allowing interactions again.
- The system automatically validates block status when attempting to interact with other users (e.g., sending friend requests).

## Technical design

The feature relies on a `blocked_users` map (type `map[string]struct{}`) maintained within the `client` structure in `server/client.go`. 
The `server/client.go` file contains the core logic for managing the `blocked_users` map, including thread-safe access protection using `sync.RWMutex`.

Command routing is handled in `client/message.go`, which maps the user's input commands to the corresponding server-side methods:
- `/user/block` routes to logic that checks if the target is already blocked, prevents self-blocking, and updates the `blocked_users` map.
- `/user/unblock` routes to logic that checks if the target is blocked, and if so, deletes the entry from the `blocked_users` map.

## Entry points

- **CLI Commands:**
  - `client/message.go:397` (Command routing for `/user/block`)
  - `client/message.go:406` (Command routing for `/user/unblock`)

## Dependencies

- `server/client.go`: Contains the `client` struct with the `blocked_users` state.
- `client/message.go`: Contains the command routing logic.

## Open questions

- None.

## Revision

- Initial draft: scope and behaviour from implementation reads in `server/client.go` and `client/message.go`.

## Evidence index

- `server/client.go:325-342`: Implementation of user blocking logic.
- `server/client.go:344-358`: Implementation of user unblocking logic.
- `server/client.go:367-391`: Implementation of check block status logic.
- `client/message.go:397-405`: Command routing for `/user/block`.
- `client/message.go:406-414`: Command routing for `/user/unblock`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Block a user](./FE-05_UC-01.md)
- [Unblock a user](./FE-05_UC-02.md)
- [Check block status](./FE-05_UC-03.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



