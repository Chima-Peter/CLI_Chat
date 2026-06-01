# Entity DM-03 — Message

**Slug:** message · **File:** DM-03-message.md

## Summary

The `Message` struct is the primary data transfer object (DTO) used for communication between the client and server. It acts as the backbone of the application's JSON wire protocol.

## Purpose and lifecycle

- **Purpose:** Encapsulates actions, responses, and associated data payloads for network communication.
- **Lifecycle:**
    1. **Construction:** Created by the client (for requests) or server (for responses/notifications).
    2. **Serialization:** Marshaled into JSON format for transmission.
    3. **Transmission:** Sent over network connections (TCP sockets).
    4. **Deserialization:** Unmarshaled by the receiver to interpret the `ActionType` and process the `Payload`.

## Invariants

- Must be valid JSON-serializable.
- The `Action` field must correspond to a defined `ActionType`.
- Transient: Objects exist only during the communication process and are not persisted in the database.

## Storage mapping

- **Storage:** Not applicable. This entity is transient and memory-resident only, used for network serialization/deserialization.
- **Protocol:** `protocol/message.go`.

## Fields overview

- `Action` (`ActionType`): Defines the operation type.
- `ResponseMsg` (`string`): Contains informative messages for the user.
- `Payload` (`json.RawMessage`): Contains data relevant to the specific action.

## Relationships overview

- Interacts with all features of the application as the base communication unit.

## Revision

- Initial draft: Entity definition based on `protocol/message.go`.
- Fixed Evidence index to remove forbidden citations.

## Evidence index

- protocol/message.go:54-58
- server/handle_conn.go:48
- client/session.go:25

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Action](./DM-03_FD-01-action.md)
- [ResponseMsg](./DM-03_FD-02-response-msg.md)
- [Payload](./DM-03_FD-03-payload.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



