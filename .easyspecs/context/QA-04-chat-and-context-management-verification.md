# Feature QA-04 — Chat and context management verification

**Slug:** chat-and-context-management-verification  
**Output file:** QA-04-chat-and-context-management-verification.md

## Summary

This feature encompasses the verification of the chat functionality and the associated context management within the application. It ensures that users can correctly switch between different interaction contexts—specifically, talking in a room or directly messaging a friend—and that messages are routed accurately based on the client's current context.

## Scope

- **In scope**:
    - Context switching mechanisms (e.g., entering a room, initiating a DM).
    - Message routing and delivery based on the established `current_context` of a client.
    - Verification of message delivery within rooms and as direct messages (DM).
- **Out of scope**:
    - Creation, deletion, or modification of rooms (covered by `QA-02`).
    - Adding, removing, or listing friends (covered by `QA-03`).
    - Initial user authentication and login flow (covered by `QA-01`).

## Functional behaviour

1.  **Context Switching**: Users initiate a context switch to a specific room or a specific friend.
2.  **Message Sending**: Once a context is established, the user sends a message.
3.  **Routing**: The server, upon receiving a message, checks the `current_context` of the authenticated client.
4.  **Delivery**: Based on the `current_context`, the server routes the message appropriately:
    - If `current_context == "room"`, the message is broadcast to all members of `cl.room`.
    - If `current_context == "friend"`, the message is sent directly to `cl.current_friend`.

## Technical design

The client state in the system is represented by the `client` struct, which maintains critical context-related fields: `current_context`, `room`, and `current_friend`.

- **State Management**: When a user connects and authenticates, their context is initialized. As the user navigates the application, the `server` updates these fields to reflect their active focus.
- **Routing Logic**: The `dispatchMessage` method in the server handles incoming messages. It examines the `current_context` to determine the target for message delivery. This routing ensures messages do not get misdirected when a user is in different conversation contexts.

## Entry points

- `server.HandleConn(net.Conn)`: Entry point for client connections, initiating the message loop.
- `server.dispatchMessage(cl *client, req *Message)`: Core handler for incoming requests, including message routing based on `current_context`.

## Dependencies

- `server/client.go`: Defines the `client` struct and the `current_context` field.
- `server/server.go`: Contains logic for updating client context and routing messages.
- `server/handle_conn.go`: Contains the dispatching mechanism that orchestrates the message routing process.

## Open questions

None.

## Revision

- Initial draft: scope, behavior, and technical design grounded in implementation details.

## Evidence index

- `server/handle_conn.go:23`: Initialization of `current_context` to `contextNone` in the `client` struct.
- `server/client.go:23`: Definition of `current_context` string field in the `client` struct.
- `server/server.go:328`: Updating `cl.current_context` to `contextRoom`.
- `server/server.go:354`: Updating `cl.current_context` to `contextFriend`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Switch context successfully](./QA-04_TC-01-tc01.md)
- [Send message to room successfully](./QA-04_TC-02-tc02.md)
- [Send message to current context successfully](./QA-04_TC-03-tc03.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



