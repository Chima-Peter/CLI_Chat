# Feature FE-01 — User authentication

**Slug:** user-authentication  
**Output file:** FE-01-user-authentication.md

## Summary

The User authentication feature handles the initial establishment of a client's identity within the system. It enforces that all clients provide a unique nickname before accessing other functional capabilities, thereby establishing an authenticated session.

## Scope

- **In scope:**
  - Client connection initialisation.
  - Nickname prompt and submission.
  - Nickname validation (non-empty, case-insensitive check against "anonymous").
  - Uniqueness enforcement (checking existing clients).
  - Transition of client state to "authenticated".
- **Out of scope:**
  - Permanent user accounts or password-protected authentication (not yet implemented).
  - User signup (explicitly marked as unimplemented).

## Functional behaviour

When a client connects to the server, they are initially in an unauthenticated state. The system prompts the client for a nickname. The client must submit a nickname, which is then validated for uniqueness across all currently connected clients. Once a valid, unique nickname is accepted, the client's `authenticated` flag is set to `true`, and they gain access to the main message dispatching functionality. If a client attempts to use other features before authentication, the system prompts them to log in first.

## Technical design

The authentication logic is primarily handled within the `server/handle_conn.go` file. The `HandleConn` method serves as the entry point for new connections, initializing the client state and prompting for a login.

Incoming messages are processed through `dispatchMessage`, which checks the `authenticated` status of the client before allowing access to restricted actions. The `handleLogin` method manages the nickname submission, validation (ensuring it is not empty and not "anonymous"), and uniqueness check using `s.isNickTaken`.

## Entry points

- **`server.HandleConn`**: The main entry point for new network connections, initiating the authentication workflow.
- **`server.handleLogin`**: The handler for processing nickname submissions.

## Dependencies

- **`server/utils.go`**: Provides the `isNickTaken` method used for uniqueness enforcement.
- **`net` package**: For socket connection handling.
- **`encoding/json`**: For processing JSON-encoded messages.

## Open questions

None.

## Revision

- Initial draft: scope, behaviour, and design based on implementation in `server/handle_conn.go` and `server/utils.go`.

## Evidence index

- `server/handle_conn.go:14-56`: `HandleConn` entry point for new connections.
- `server/handle_conn.go:60-71`: `dispatchMessage` authentication check.
- `server/handle_conn.go:164-192`: `handleLogin` implementation, including validation and uniqueness checks.
- `server/utils.go:116-123`: `isNickTaken` helper method for uniqueness enforcement.

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



