# Feature FE-01 — Authentication

**Slug:** authentication  
**Output file:** FE-01-authentication.md

## Summary

The Authentication feature handles user connection, nickname registration, and session management. It ensures that only users with a unique, valid nickname can access the system's functional features.

## Scope

- **In scope:**
  - User connection handling.
  - Login via nickname.
  - Validation of nicknames (non-empty, not "anonymous", unique).
  - Session authentication state management.
  - Logout functionality.
- **Out of scope:**
  - Robust user account persistence (signup is currently a stub).
  - Password-based authentication for user accounts (only nickname-based).

## Functional behaviour

When a client connects, they are prompted to enter a nickname. The server validates the nickname:
- It must not be empty.
- It cannot be "anonymous".
- It must not be taken by another connected user.

Once successfully logged in, the user's connection is authenticated, enabling access to other features like messaging, room management, and friend management. If a client attempts to use these features without being authenticated, the server denies the request.

## Technical design

Authentication is managed in `server/handle_conn.go`. The server maintains a `client` struct (defined in `server/handle_conn.go`) which includes an `authenticated` boolean field. 

The `HandleConn` method serves as the entry point for new connections, initializing the client state and looping to process messages. The `dispatchMessage` method acts as a gatekeeper, checking the `authenticated` flag before processing actions. If unauthenticated, only `LOGIN` and `LOGOUT` actions are accepted.

The `handleLogin` method performs the validation logic and updates the client's `nick` and `authenticated` status upon success.

## Entry points

- `server/handle_conn.go`: `HandleConn(conn net.Conn)` - Initial entry point for new connections.
- `server/handle_conn.go`: `handleLogin(cl *client, nick string)` - Handles the login action.

## Dependencies

- `github.com/google/uuid`: Used for generating unique client IDs upon connection.

## Open questions

None.

## Revision

- Initial draft: scope and behaviour from implementation reads.

## Evidence index

- `server/handle_conn.go:14-56`: Client connection handling and message processing loop.
- `server/handle_conn.go:60-70`: Authentication gatekeeper logic in `dispatchMessage`.
- `server/handle_conn.go:164-192`: `handleLogin` logic for nickname validation and authentication status update.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Log in with a valid nickname](./FE-01_UC-01.md)
- [Attempt to log in with an invalid or taken nickname](./FE-01_UC-02.md)
- [Log out](./FE-01_UC-03.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



