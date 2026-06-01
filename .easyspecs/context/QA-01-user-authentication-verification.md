# Feature QA-01 — User authentication verification

**Slug:** user-authentication-verification
**Output file:** QA-01-user-authentication-verification.md

## Summary
The "User authentication verification" feature is responsible for handling the user login process, ensuring that each client is identified by a unique, valid nickname. This is a critical security and identity step in the application.

## Scope
- **In-scope:** Client authentication flow, nickname validation (non-empty, non-"anonymous"), and uniqueness check.
- **Out-of-scope:** Password authentication (not implemented), social logins, and persistent account management.

## Functional behaviour
When a user connects or attempts to login, the system prompts them for a nickname. The following validations are performed:
1. **Already Authenticated:** If a user is already logged in, they are informed and asked to log out first.
2. **Empty Nickname:** Nicknames cannot be empty or solely whitespace.
3. **Reserved Name:** The nickname "anonymous" (case-insensitive) is disallowed.
4. **Uniqueness:** The server verifies that the chosen nickname is not already in use by another client.

If all validations pass, the client is marked as authenticated and their nickname is set.

## Technical design
The authentication logic is centralized in the `handleLogin` function within the `server/handle_conn.go` module. The server maintains client state (authenticated status and nickname) under a mutex lock to ensure thread safety during the update.

## Entry points
- `server/handle_conn.go`: `handleLogin` function. This is triggered when a client sends a login request.

## Dependencies
- `server/utils.go`: `isNickTaken` function, which iterates through active clients to enforce nickname uniqueness.

## Open questions
None.

## Revision
- Initial draft: scope, behaviour, and technical design from implementation analysis.

## Evidence index
- `server/handle_conn.go:164-192`: `handleLogin` implementation detailing the authentication flow and validation rules.
- `server/utils.go:116`: `isNickTaken` implementation used for enforcing nickname uniqueness.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Successful login with a valid nickname](./QA-01_TC-01-successful-login.md)
- [Failed login with empty nickname](./QA-01_TC-02-failed-login-empty-nickname.md)
- [Failed login with 'anonymous' nickname](./QA-01_TC-03-failed-login-anonymous.md)
- [Failed login with already taken nickname](./QA-01_TC-04-failed-login-taken-nickname.md)
- [Logout after successful login](./QA-01_TC-05-logout-successful.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



