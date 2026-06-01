# Use case TC-05 — Logout after successful login (Feature QA-01)

## Summary

A successfully authenticated client requests to log out of the system. The server processes this request by sending a confirmation message and closing the connection.

## Actors and stakeholders

- **Client**: The actor initiating the logout request.
- **Server**: The system processing the request and managing client connectivity.

## Preconditions

- The client must be successfully authenticated (i.e., `client.authenticated` is `true`).

## Data inputs and validation

- **Input**: A JSON `Message` object where `Action` is set to `LOGOUT`.
- **Validation**:
  - The `dispatchMessage` function checks if the client is authenticated (`if !cl.authenticated` in `server/handle_conn.go:60`).
  - If authenticated, it routes to the `LOGOUT` case in the switch statement (`server/handle_conn.go:73,78`).

## Main flow (user- or operator-visible)

1. **Client** sends a `LOGOUT` action message to the **Server**.
2. **Server** receives the message, identifies it as a logout request.
3. **Server** sends a `DONE` response message with "Logged out." to the **Client**.
4. **Server** returns `true` from `dispatchMessage` to indicate the connection should be closed.
5. **Server** closes the client connection.

## Code flow

1. **Entrypoint**: `HandleConn` in `server/handle_conn.go` receives the message via `decoder.Decode`.
2. **Dispatch**: `s.dispatchMessage(cl, &request)` is called.
3. **Auth Check**: `dispatchMessage` verifies `cl.authenticated`.
4. **Logic**: `LOGOUT` case is triggered:
   - `cl.send_user_message(...)` sends "Logged out.".
   - Returns `true` to signal connection closure.
5. **Connection Closure**: `HandleConn` loop exits, and `defer conn.Close()` (and `defer s.LogUserOut(cl)`) handles the cleanup.

### Code flow (Implementation)

```mermaid
%%{init: {'theme':'neutral'}}%%
flowchart TD
    A[Client] -->|{"Action": "LOGOUT"}| B[Server.dispatchMessage]
    B --> C{Authenticated?}
    C -->|Yes| D[Case LOGOUT]
    D --> E[cl.send_user_message: "Logged out."]
    D --> F[Return true]
    F --> G[HandleConn loop terminates]
    G --> H[defer conn.Close()]
```

## Alternate flows

- **Logout before login**: If a client attempts to log out before being authenticated, the `!cl.authenticated` block handles it:
  - `dispatchMessage` receives the `LOGOUT` action (`server/handle_conn.go:65-66`).
  - Returns `true` to close the connection immediately without extra messages.

## Postconditions

- Client connection is closed.
- Client is removed from server's active clients list (via `LogUserOut`).

## Errors and edge cases

- N/A (The logout flow is simple and robust).

## Technical mapping

- **Message handling**: `server/handle_conn.go:52`
- **Logout logic**: `server/handle_conn.go:78-80`
- **Connection cleanup**: `server/handle_conn.go:15, 41`

## Related scenarios

- None.

## Revision

- Initial draft: Defined flow, actors, preconditions, data inputs, code flow, and evidence.

## Evidence index

- `server/handle_conn.go:14-15` — Entrypoint (HandleConn) and connection closure defer.
- `server/handle_conn.go:52` — Dispatcher entrypoint.
- `server/handle_conn.go:78-80` — Core business logic for logout.
- `server/handle_conn.go:41` — Cleanup logic (`LogUserOut`).

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [User authentication verification](./QA-01-user-authentication-verification.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



