# Feature QA-05 — User management verification

**Slug:** user-management-verification  
**Output file:** QA-05-user-management-verification.md

## Summary

The User management feature enables users to control their interactions by blocking or unblocking other users, and checking the status of other users (online/offline status, current room, and blocked status).

## Scope

This feature covers the functionality of user-to-user blocking, unblocking, and status retrieval. It does not cover authentication, room creation, or friend request handling, which are part of other features (FE-01, FE-02, FE-03).

## Functional behaviour

Users can interact with other users via the following commands:
- **Block:** Prevents the target user from initiating interactions (e.g., friend requests, messaging).
- **Unblock:** Reverses the block action, allowing interactions again.
- **Get Status:** Retrieves information about a user, including:
  - Online/offline status.
  - Current room (if applicable).
  - Whether the target user is blocked by the requester.

## Technical design

User information and management state are maintained within the `client` structure in `server/client.go`. The server handles these requests by resolving the target user and delegating the operation to the client's methods.

- **Blocking/Unblocking:** Modifies the `blocked_users` map in the `client` structure (`server/client.go:29`, `server/client.go:323`, `server/client.go:344`).
- **Status Retrieval:** Queries the `online`, `room`, and `blocked_users` fields of the target user (`server/client.go:366`).

## Entry points

The following methods in `server/server.go` serve as entry points for user management requests:
- `server.BlockUser` (`server/server.go:651`)
- `server.UnblockUser` (`server/server.go:660`)
- `server.GetUserStatus` (`server/server.go:669`)

## Dependencies

- `server/client.go` (Client state and method definitions).

## Open questions

None.

## Revision

- Initial draft: scope, behaviour, and technical design grounded in `server/server.go` and `server/client.go`. Corrected the source reference range for the implementation.

## Evidence index

- Implementation of user management handlers: `server/server.go:651-676`
- Implementation of client-side user management methods: `server/client.go:323-394`

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Block user](./QA-05_TC-01-tc01.md)
- [Unblock user](./QA-05_TC-02-tc02.md)
- [Get user status](./QA-05_TC-03-tc03.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



