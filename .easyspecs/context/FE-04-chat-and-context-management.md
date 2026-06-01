# Feature FE-04 — Chat and context management

**Slug:** chat-and-context-management  
**Output file:** FE-04-chat-and-context-management.md

## Summary

The "Chat and context management" feature enables users to navigate conversations within different contexts (rooms or direct messages with friends) and exchange messages. It supports switching between active conversation contexts, sending messages to a room, and sending direct messages to friends.

## Scope

- **In scope:**
  - Setting and switching the active context (room or friend).
  - Sending messages to the currently active context (room or friend).
  - Explicitly sending messages to a specific room or friend.
- **Out of scope:**
  - File transfers (handled in a separate feature).
  - Room management (handled in FE-02).
  - Friend management (handled in FE-03).

## Functional behaviour

Users can set their active context using the `SWITCH_CONTEXT` command. Once a context is active, messages sent via `SEND_MSG` are automatically routed to that context. Alternatively, users can use `MESSAGE_ROOM` or `MESSAGE_FRIEND` to send messages without needing to switch their active context.

## Technical design

The feature relies on the message dispatch loop in `server/handle_conn.go`, which routes incoming `Message` payloads to appropriate server-side handlers based on the `Action` field.

- **Message Dispatching:** `server/handle_conn.go` functions as the main entry point for processing chat actions.
- **Context Management:** `server/server.go` manages the state for context switching (`SwitchContext`) and dispatching messages based on the client's current context (`SendContextMessage`).

## Entry points

- **Actions (defined in `server/handle_conn.go`):**
  - `SWITCH_CONTEXT`: `server/handle_conn.go:116`
  - `MESSAGE_ROOM`: `server/handle_conn.go:126`
  - `SEND_MSG`: `server/handle_conn.go:128`
  - `MESSAGE_FRIEND`: `server/handle_conn.go:138`
- **Handlers:**
  - `SwitchContext`: `server/server.go:301`
  - `SendContextMessage`: `server/server.go:369`

## Dependencies

- **FE-02 (Room management):** Dependent for room resolution.
- **FE-03 (Friend management):** Dependent for friend resolution.

## Open questions

None.

## Revision

- Initial draft: scope, behaviour, design, and entry points derived from `server/handle_conn.go` and `server/server.go`.

## Evidence index

- Message dispatching loop and action mapping: `server/handle_conn.go:58-158`
- Context switching handler definition: `server/server.go:301`
- Context message sending handler definition: `server/server.go:369`

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



