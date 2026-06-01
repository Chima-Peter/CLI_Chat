# Use case TC-03 — Failed login with 'anonymous' nickname (Feature QA-01)

## Summary

The system rejects any login attempt where the provided nickname is "anonymous" (case-insensitive), prompting the user to provide a different, valid nickname.

## Actors and stakeholders

- **User**: Attempts to connect to the server and authenticate with an invalid nickname.
- **Server**: Validates the incoming nickname against the "anonymous" restriction.

## Preconditions

- The client is connected to the server.
- The client is not yet authenticated.

## Data inputs and validation

- **Input**: Nickname provided during the login process (`LOGIN` action).
- **Validation**:
  - The nickname must not be "anonymous" (case-insensitive check using `strings.EqualFold`).
  - The validation runs in `server/handle_conn.go` within `handleLogin` method.
- **Failure behaviour**:
  - The server sends a `LOGIN` message back to the client with the message: "Nickname cannot be anonymous. Enter a valid nickname: ".
  - The connection remains open, and the client remains in an unauthenticated state, allowing retry.

## Main flow (user- or operator-visible)

1. The user connects to the server and receives an initial "Enter your nickname: " prompt.
2. The user sends a `LOGIN` message with the payload "anonymous" (or "Anonymous", "ANONYMOUS", etc.).
3. The server validates the nickname and detects it is prohibited.
4. The server sends a `LOGIN` response to the client prompting for a valid nickname.

## Code flow

### Request path (implementation)

1. **Entrypoint**: `server/HandleConn` initiates the connection loop.
2. **Action dispatch**: `server/dispatchMessage` receives the `LOGIN` action.
3. **Logic**: `server/handleLogin` is called with the nickname.
4. **Validation**: `server/handleLogin` checks `strings.EqualFold(nick, "anonymous")`.
5. **Rejection**: If invalid, `cl.promptLogin(...)` is called to send the error message, and the function returns, keeping `authenticated` as `false`.

### Mermaid

```mermaid
%%{init: {'theme':'neutral'}}%%
flowchart TD
    A[Client] -->|Sends LOGIN action| B[server.dispatchMessage]
    B -->|Calls| C[server.handleLogin]
    C -->|Check: strings.EqualFold| D{Nickname == "anonymous"?}
    D -- Yes --> E[cl.promptLogin]
    E -->|Sends ERROR response| A
    D -- No --> F[Proceed to authentication]
```

## Alternate flows

- **Valid nickname**: The nickname is not "anonymous", not empty, and not taken; the login proceeds to completion.
- **Empty nickname**: `TC-02` (Handled before the anonymous check).
- **Taken nickname**: `TC-04` (Handled after the anonymous check).

## Postconditions

- The user remains unauthenticated.
- The client connection is maintained to allow for a retry.

## Errors and edge cases

- **Case sensitivity**: "Anonymous", "ANONYMOUS", and "anonymous" are all rejected.

## Technical mapping

- The validation relies on `strings.EqualFold` for case-insensitive comparison.

## Related scenarios

- None.

## Revision

- Initial draft: data inputs, code flow, and Evidence tied to handlers and services.

## Evidence index

- `server/handle_conn.go:14-56` — Entry point (`HandleConn`)
- `server/handle_conn.go:61-64` — Dispatcher for `LOGIN` action
- `server/handle_conn.go:174-177` — Validation logic for anonymous nickname rejection
- `server/handle_conn.go:160-162` — Prompt helper for login retries

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [User authentication verification](./QA-01-user-authentication-verification.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



