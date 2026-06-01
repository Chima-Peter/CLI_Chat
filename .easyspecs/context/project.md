# CLI Chat — project

## Summary

CLI Chat is a terminal-based chat application written in Go. It implements a TCP-based chat system that supports rooms, friend requests, direct messaging, and server-driven prompts. The system uses a JSON-based wire protocol to facilitate communication between an interactive CLI client and a concurrent TCP server.

## Overall goal

The goal of this project is to provide a functional, lightweight, terminal-based chat platform that enforces structured communication (via JSON protocol) while maintaining user state, room membership, and friend relationships in memory. Success is defined by the ability to reliably connect, authenticate, and communicate across multiple rooms and friend contexts.

## Illustrative capabilities and use cases

> **Note:** The lists below are **illustrative only**, not an exhaustive inventory of every feature or workflow.

### Representative functionalities

- **Concurrent TCP Server**: Handles multiple client connections using goroutines.
- **JSON Protocol**: Structured, line-delimited communication for commands and messages.
- **Room Management**: Create, join, leave, edit, and invite functionalities.
- **Friend Management**: Request, accept, and remove friends for direct messaging.
- **Context Switching**: `/switch` command allows users to toggle plain-text destinations between rooms and individual friends.
- **Login/Auth**: Mandatory nickname prompt and uniqueness enforcement on connection.

### Example use cases

- **Collaborating in a room**: A user connects, creates a room, joins it, switches context to that room, and starts broadcasting messages.
- **Direct messaging**: A user adds a friend, and once the request is accepted, they can DM directly using `/chat/dm` or switch context to that friend.

## Purpose

The application serves as a chat platform for terminal users. It maintains in-memory state for active sessions, room memberships, and friend lists. The client and server communicate via a JSON protocol over TCP, where slash commands are parsed by the client and sent as structured actions.

## Repository role

This repository contains the source code for both the CLI client and the TCP server. It is a standalone application intended for local execution and experimentation with Go-based networked systems.

## Tech stack

| Area | Choice |
| ---- | ------ |
| Language | Go (1.26.3+) |
| Protocol | TCP with JSON payloads |
| Data | In-memory (no persistence) |

## Scripts

| Command | Description |
| ------- | ----------- |
| `go build -o server ./cmd/server` | Builds the server binary |
| `go build -o client ./cmd/client` | Builds the client binary |
| `./server` | Starts the server |
| `./client` | Starts the client |

## Source layout (high level)

| Path | Role |
| ---- | ---- |
| `cmd/` | Entry points for server and client |
| `protocol/` | Shared JSON protocol definitions |
| `client/` | Client-side logic (UI, command parsing, session management) |
| `server/` | Server-side logic (handlers, room management, state) |

## Configuration and environment

- Server listens on port `:8888`.
- Requires Go 1.26.3+ for building.

## Related documentation

- `README.md` — Main documentation and quick start guide.
- `architecture.md` (Expected in the same directory: `.easyspecs/context/`) — Deeper architectural details and system flow diagrams.

## Out of scope (current phase)

- Data persistence (restart resets all state).
- `SIGN_UP` action implementation.
- `SEND_FILE` functionality.

## Revision

- Initial draft: summary, goal, and tech stack from README.

## Evidence index

- `go.mod:3` (Go version prerequisite)
- `protocol/message.go:5-58` (JSON protocol definition and message structure)
- `README.md:1-245` (CLI Chat overview, features, and usage)

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

### Architecture narrative

- [Application architecture](./architecture.md)

### Features
- [Room management](./FE-02-room-management.md)
- [Friend management](./FE-03-friend-management.md)
- [Chat and context management](./FE-04-chat-and-context-management.md)
- [User management](./FE-05-user-management.md)
- [User authentication](./FE-01-user-authentication.md)
### Experiences
*None.*
### Services
- [Server](./SV-01-server.md)
- [Client](./SV-02-client.md)
### Tech stack
- [Go](./TS-01-go.md)
- [github.com/chzyer/readline](./TS-02-readline.md)
- [github.com/google/uuid](./TS-03-uuid.md)
- [golang.org/x/sys](./TS-04-x-sys.md)
### Infrastructure
- [CLI Client](./IN-01-cli-client.md)
- [CLI Server](./IN-02-cli-server.md)
- [TLS Configuration](./IN-03-tls-configuration.md)
- [Dependency Management](./IN-04-dependency-management.md)
### QA
- [User authentication verification](./QA-01-user-authentication-verification.md)
- [Room management verification](./QA-02-room-management-verification.md)
- [Friend management verification](./QA-03-friend-management-verification.md)
- [Chat and context management verification](./QA-04-chat-and-context-management-verification.md)
- [User management verification](./QA-05-user-management-verification.md)
### Data model
- [Room](./DM-01-room.md)
- [Client](./DM-02-client.md)
- [Message](./DM-03-message.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

*None.*
<!-- easyspecs-nav:parents:end -->



