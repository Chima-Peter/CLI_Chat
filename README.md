# CLI Chat

A command-line chat application in Go with a TCP server and an interactive JSON protocol client.

## Features

- **Multi-client TCP server** — concurrent connections via goroutines
- **JSON wire protocol** — newline-delimited `protocol.Message` objects
- **Chat rooms** — create, join, leave, edit, passwords, invites, member management
- **Friends & DMs** — friend requests, direct messages, block/unblock, online status
- **Server-driven prompts** — password and other flows guided by server actions
- **Rich client commands** — slash commands map to protocol actions (`/help` lists all)

## Project Structure

```
CLI_Chat/
├── cmd/
│   ├── client/          # Client entry point
│   └── server/          # Server entry point
├── protocol/
│   └── message.go       # Shared action types and Message struct
├── client/
│   ├── client.go        # Connect and bootstrap
│   ├── session.go       # Readline session, server message handling
│   ├── message.go       # Build commands and protocol messages
│   └── help.go          # /help text and startup banner
├── server/
│   ├── server.go        # Room and friend handlers
│   ├── handle_conn.go   # Connection loop and dispatch
│   ├── client.go        # Per-client state and messaging
│   ├── rooms.go         # Room logic
│   └── utils.go         # Payload parsing and lookups
└── go.mod
```

## Installation

### Prerequisites

- Go 1.26.3 or higher

### Build

```bash
git clone https://github.com/chima/CLI_Chat.git
cd CLI_Chat

go build -o server ./cmd/server
go build -o client ./cmd/client
```

## Usage

### Start the server

```bash
./server
```

Listens on `localhost:8888`.

### Start the client

```bash
./client
```

Use `/help` in the client for the full command list. Quick start:

```
/auth/login alice
/room/create general
/room/join general
Hello everyone!
/room/list
/invite/mine
/auth/logout
```

## Protocol

Each message on the wire is one JSON object per line:

```json
{
  "action": 12,
  "response_msg": "Welcome to room.",
  "payload": {"room": "general", "room_id": "..."}
}
```

| Category | Actions |
|----------|---------|
| Auth | `SIGN_UP`, `LOGIN`, `LOGOUT` |
| Rooms | `CREATE_ROOM`, `SET_ROOM_PASSWORD`, `JOIN_ROOM`, `LEAVE_ROOM`, `DELETE_ROOM`, `EDIT_ROOM`, `GET_ROOM_PASSWORD`, `DELETE_MEMBER`, invites, `GET_ROOM_MEMBERS`, `LIST_ROOMS`, `LIST_MY_ROOMS`, `LIST_MY_ROOM_INVITES` |
| Chat | `SEND_MSG`, `SEND_FILE` |
| Friends | `SEND_FRIEND_REQUEST`, `ACCEPT_FRIEND_REQUEST`, `MESSAGE_FRIEND`, `GET_FRIENDS`, `SEE_FRIEND_REQUEST`, `DELETE_FRIEND`, `BLOCK_USER`, `UNBLOCK_USER`, `GET_USER_STATUS` |
| Response | `ERR`, `DONE` |

### Server-driven prompts

When the server responds with an action other than `DONE`, `ERR`, `SEND_MSG`, or `MESSAGE_FRIEND`, the client treats that action as the next reply type (e.g. `GET_ROOM_PASSWORD` after joining a private room). Your next line is sent with the same action and merged payload.

### Client commands

All slash commands are handled locally in `buildCommandMessage` and sent as protocol messages. `/help` is **client-only** and does not hit the server.

| Command | Protocol |
|---------|----------|
| **Auth** | |
| `/auth/login <name>` | `LOGIN` |
| `/auth/signup <name>` | `SIGN_UP` |
| `/auth/logout` | `LOGOUT` |
| **Room** | |
| `/room/create <name>` | `CREATE_ROOM` |
| `/room/join <name>` | `JOIN_ROOM` |
| `/room/leave [name]` | `LEAVE_ROOM` |
| `/room/delete <name>` | `DELETE_ROOM` |
| `/room/edit <room> <new> [max]` | `EDIT_ROOM` |
| `/room/members <room>` | `GET_ROOM_MEMBERS` |
| `/room/kick <room> <user>` | `DELETE_MEMBER` |
| `/room/invite <room> <user>` | `SEND_INVITE_REQUEST` |
| `/room/invites [room]` | `SEE_GROUP_INVITE_REQUEST` (owner: users you invited) |
| `/room/list` | `LIST_ROOMS` |
| `/room/mine` | `LIST_MY_ROOMS` |
| **Invite** | |
| `/invite/mine` | `LIST_MY_ROOM_INVITES` |
| `/invite/accept <room>` | `ACCEPT_GROUP_INVITE_REQUEST` |
| `/invite/decline <room>` | `DELETE_GROUP_INVITE_REQUEST` |
| **Chat** | |
| `/chat/send <text>` or plain text | `SEND_MSG` |
| `/chat/file <path>` | `SEND_FILE` |
| `/chat/dm <user> <msg>` | `MESSAGE_FRIEND` |
| **Friend** | |
| `/friend/add <user>` | `SEND_FRIEND_REQUEST` |
| `/friend/accept <user>` | `ACCEPT_FRIEND_REQUEST` |
| `/friend/remove <user>` | `DELETE_FRIEND` |
| `/friend/requests` | `SEE_FRIEND_REQUEST` |
| `/friend/list` | `GET_FRIENDS` |
| **User** | |
| `/user/block <user>` | `BLOCK_USER` |
| `/user/unblock <user>` | `UNBLOCK_USER` |
| `/user/status <user>` | `GET_USER_STATUS` |

Payload fields commonly use `room`, `room_id`, `nick` / `username`, `user_id`, `message`, `password`, `new_room`, `max_size`.

## Architecture

### Server

- Accepts TCP connections; one goroutine per `HandleConn`
- Dispatches incoming `action` values in `dispatchMessage`
- Replies with `send_user_message` (JSON encode to client)
- Rooms keyed by ID; users and friends keyed by user ID

### Client

- `readFromServer` decodes JSON in the background
- `handleServerMessage` prints responses and tracks server prompts
- `inputLoop` uses readline for input; clears partial input when a server prompt arrives
- Plain text (no `/`) sends `SEND_MSG` to the current room

## Not yet implemented

- `SIGN_UP`, `SEND_FILE` (server returns not implemented)
- Persistence, nickname uniqueness enforcement, full block enforcement on all paths

## License

MIT License

## Author

Chima

## Contributing

Issues and pull requests are welcome.
