# Feature FE-04 — Messaging and Context

**Slug:** messaging-and-context
**Output file:** FE-04-messaging-and-context.md

## Summary

The Messaging and Context feature facilitates communication within the chat application by managing user context (room or friend) and handling message routing, including plain text messages and direct file transfers.

## Scope

- **In Scope:**
    - Switching active context between rooms and friends.
    - Sending plain text messages to the active context (room or friend).
    - Managing direct messages between friends.
    - Handling file transfer initiation and metadata exchange.
- **Out of Scope:**
    - User authentication (handled by FE-01).
    - Room creation/management lifecycle (handled by FE-02).
    - Friend list management lifecycle (handled by FE-03).
    - User blocking/status (handled by FE-05).

## Functional behaviour

Users can use slash commands to switch their active communication context to either a room or a friend. Once the context is set, subsequent messages sent by the user are automatically routed to the active room or friend. File transfers require a request/response mechanism to exchange connection details between friends.

## Technical design

The feature is implemented across both the server and the client.
- The server (`server/server.go`) maintains the active context for each client, routes messages based on the current context, and facilitates file transfer port exchange.
- The client (`client/session.go`) manages the input loop, slash command parsing, and file server/client interactions to support direct file transfers.

## Entry points

- `/switch <room|friend> <name>`: Used to switch context (`server/server.go:301-367`).
- Plain text messages (when context is set): Routed based on current context (`server/server.go:369-401`).
- `/send <friend> <path>`: Initiates a file transfer (`client/session.go:223-231`).

## Dependencies

- **`protocol`**: Defines message structures and actions.
- **`github.com/chzyer/readline`**: Used by the client for input handling.

## Open questions

None

## Revision

- Initial draft: scope and behaviour from implementation reads.

## Evidence index

- `server/server.go:301-367`: `SwitchContext` implementation.
- `server/server.go:369-401`: `SendContextMessage` implementation.
- `client/session.go:193-256`: `inputLoop` implementation.
- `client/session.go:223-231`: File transfer initiation logic in `inputLoop`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Switch context](./FE-04_UC-01.md)
- [Send message to room](./FE-04_UC-02.md)
- [Send message to current context](./FE-04_UC-03.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



