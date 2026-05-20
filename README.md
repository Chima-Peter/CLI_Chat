# CLI Chat

A terminal chat application in Go: a TCP server and an interactive client that speak a JSON protocol. Commands are grouped like REST API routes (`/room/join`, `/friend/add`, `/chat/dm`), and the server tracks an **active chat context** so plain text goes to the right room or friend.

## Features

- **Concurrent TCP server** — one goroutine per connection
- **JSON wire protocol** — one `protocol.Message` per line
- **Rooms** — create, join, leave, edit, passwords, invites, members, kick
- **Friends** — requests, list, remove, direct messages
- **Context switching** — `/switch` sets whether plain text goes to a room or a friend
- **Explicit chat commands** — `/chat/room` and `/chat/dm` target a room or friend by name without changing context
- **Server-driven prompts** — e.g. room password after create/join
- **Numbered list responses** — rooms, members, invites, friends shown as `1. …`, `2. …`

## Project structure

```
CLI_Chat/
├── cmd/
│   ├── client/          # Client entry point
│   └── server/          # Server entry point (listens on :8888)
├── protocol/
│   └── message.go       # Shared ActionType enum and Message struct
├── client/
│   ├── client.go        # Connect and bootstrap
│   ├── session.go       # Readline, server message loop, prompts
│   ├── message.go       # Slash commands → protocol messages
│   └── help.go          # Grouped /help and startup banner
├── server/
│   ├── server.go        # Room, invite, context, and friend handlers
│   ├── handle_conn.go   # Per-connection read loop and dispatch
│   ├── client.go        # Per-client state (context, friends, rooms)
│   ├── rooms.go         # Room membership, broadcast, invites
│   └── utils.go         # Payload parsing and lookups
└── go.mod
```

## Installation

**Prerequisites:** Go 1.26.3+

```bash
git clone https://github.com/chima/CLI_Chat.git
cd CLI_Chat

go build -o server ./cmd/server
go build -o client ./cmd/client
```

## Usage

**Terminal 1 — server:**

```bash
./server
```

**Terminal 2 — client:**

```bash
./client
```

Type `/help` in the client for the full command list.

### Quick start

```
/auth/login alice
/room/create general
/room/join general
/switch room general
Hello everyone!
/auth/logout
```

Joining a room does **not** set your chat context. Use `/switch room <name>` before sending plain text, or use `/chat/room <name> <message>` for a one-off room message.

### DM example

```
/friend/add bob
# (bob accepts with /friend/accept alice)
/switch friend bob
Hey Bob!
```

Or without switching: `/chat/dm bob Hey Bob!`

---

## Client commands

Slash commands are parsed in `client/message.go` and sent as protocol messages. `/help` is client-only.

### Auth

| Command | Protocol action |
|---------|-----------------|
| `/auth/login <name>` | `LOGIN` |
| `/auth/signup <name>` | `SIGN_UP` (stub) |
| `/auth/logout` | `LOGOUT` |

### Room

| Command | Protocol action |
|---------|-----------------|
| `/room/create <name>` | `CREATE_ROOM` |
| `/room/join <name>` | `JOIN_ROOM` |
| `/room/leave [name]` | `LEAVE_ROOM` |
| `/room/delete <name>` | `DELETE_ROOM` |
| `/room/edit <room> <new_name> [max]` | `EDIT_ROOM` |
| `/room/members <room>` | `GET_ROOM_MEMBERS` |
| `/room/kick <room> <user>` | `DELETE_MEMBER` |
| `/room/invite <room> <user>` | `SEND_INVITE_REQUEST` |
| `/room/invites [room]` | `SEE_GROUP_INVITE_REQUEST` (owner: users you invited) |
| `/room/list` | `LIST_ROOMS` |
| `/room/mine` | `LIST_MY_ROOMS` |

### Invite

| Command | Protocol action |
|---------|-----------------|
| `/invite/mine` | `LIST_MY_ROOM_INVITES` |
| `/invite/accept <room>` | `ACCEPT_GROUP_INVITE_REQUEST` |
| `/invite/decline <room>` | `DELETE_GROUP_INVITE_REQUEST` |

### Chat

| Command | Protocol action | Notes |
|---------|-----------------|-------|
| `/switch room <name>` | `SWITCH_CONTEXT` | Active context = room |
| `/switch friend <name>` | `SWITCH_CONTEXT` | Active context = friend |
| *plain text* | `SEND_MSG` | Sends to active context only |
| `/chat/room <room> <message>` | `MESSAGE_ROOM` | Direct room message; context unchanged |
| `/chat/dm <friend> <message>` | `MESSAGE_FRIEND` | Direct DM; context unchanged |
| `/chat/file <path>` | `SEND_FILE` (stub) | Not implemented on server |

### Friend

| Command | Protocol action |
|---------|-----------------|
| `/friend/add <user>` | `SEND_FRIEND_REQUEST` |
| `/friend/accept <user>` | `ACCEPT_FRIEND_REQUEST` |
| `/friend/remove <user>` | `DELETE_FRIEND` |
| `/friend/requests` | `SEE_FRIEND_REQUEST` |
| `/friend/list` | `GET_FRIENDS` |

### User

| Command | Protocol action |
|---------|-----------------|
| `/user/block <user>` | `BLOCK_USER` |
| `/user/unblock <user>` | `UNBLOCK_USER` |
| `/user/status <user>` | `GET_USER_STATUS` |

---

## How messaging works

Three separate paths:

```
┌─────────────────────────────────────────────────────────────────┐
│  Plain text (no /)                                               │
│  → SEND_MSG { "message": "..." }                                 │
│  → Server: SendContextMessage                                    │
│     • current_context == "room"  → broadcast to cl.room          │
│     • current_context == "friend" → DM to cl.current_friend    │
│     • no context → error (use /switch first)                     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│  /chat/room <room> <message>                                     │
│  → MESSAGE_ROOM { "room": "...", "message": "..." }               │
│  → Server: SendRoomMessage (membership check, then broadcast)    │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│  /chat/dm <friend> <message>                                     │
│  → MESSAGE_FRIEND { "nick": "...", "message": "..." }            │
│  → Server: MessageFriend (friend check, online, deliver)         │
└─────────────────────────────────────────────────────────────────┘
```

### Server context state

Each connected client stores:

| Field | Purpose |
|-------|---------|
| `current_context` | `""`, `"room"`, or `"friend"` |
| `room` | Active room when context is `room` |
| `current_friend` | Active friend client when context is `friend` |

`/switch` updates these fields. Room membership (`my_rooms`) is separate: you can be in many rooms but only one **active** target for plain text at a time.

### Incoming messages

- **Room broadcast** from other members → server sends `MESSAGE_ROOM` with `response_msg` (display-only on client).
- **DM** → server sends `MESSAGE_FRIEND` to the recipient.

---

## Protocol

One JSON object per line on the wire:

```json
{
  "action": 30,
  "response_msg": "Switched to room: general",
  "payload": {"context": "room", "room": "general"}
}
```

### Action categories

| Category | Actions |
|----------|---------|
| **Auth** | `SIGN_UP`, `LOGIN`, `LOGOUT` |
| **Room** | `CREATE_ROOM`, `SET_ROOM_PASSWORD`, `JOIN_ROOM`, `LEAVE_ROOM`, `DELETE_ROOM`, `EDIT_ROOM`, `GET_ROOM_PASSWORD`, `DELETE_MEMBER`, `SEND_INVITE_REQUEST`, `SEE_GROUP_INVITE_REQUEST`, `ACCEPT_GROUP_INVITE_REQUEST`, `DELETE_GROUP_INVITE_REQUEST`, `GET_ROOM_MEMBERS`, `LIST_ROOMS`, `LIST_MY_ROOMS`, `LIST_MY_ROOM_INVITES` |
| **Chat** | `SWITCH_CONTEXT`, `MESSAGE_ROOM`, `SEND_MSG`, `SEND_FILE` |
| **Friend** | `SEND_FRIEND_REQUEST`, `ACCEPT_FRIEND_REQUEST`, `MESSAGE_FRIEND`, `GET_FRIENDS`, `SEE_FRIEND_REQUEST`, `DELETE_FRIEND` |
| **User** | `BLOCK_USER`, `UNBLOCK_USER`, `GET_USER_STATUS` |
| **Response** | `ERR`, `DONE` |

`ActionType` values are defined in order in `protocol/message.go` (iota). Client and server must be built from the same revision so action IDs match.

### Common payload fields

| Field | Used for |
|-------|----------|
| `username`, `nick` | Login, friend/room user targets |
| `room`, `room_id`, `room_name` | Room operations |
| `new_room`, `max_size` | Edit room |
| `password` | Private room create/join prompts |
| `message` | Chat body |
| `context` | `SWITCH_CONTEXT`: `"room"` or `"friend"` |
| `friend_name` | `SWITCH_CONTEXT` when context is friend |

### Server-driven prompts

If the server replies with an action other than `DONE`, `ERR`, `MESSAGE_ROOM`, `SEND_MSG`, or `MESSAGE_FRIEND`, the client treats it as a **prompt** and sends your next line with the same action and merged payload (e.g. `GET_ROOM_PASSWORD` after joining a private room).

| Prompt action | When |
|---------------|------|
| `SET_ROOM_PASSWORD` | After `/room/create` — set password (empty = public) |
| `GET_ROOM_PASSWORD` | After `/room/join` on a private room |

---

## Architecture

### Server

- `HandleConn` accepts a connection, registers a `client`, and loops on decoded `Message` values.
- `dispatchMessage` switches on `action` and calls handlers in `server.go` / `rooms.go` / `client.go`.
- Replies use `send_user_message` (JSON-encoded `Message` on the TCP connection).
- Rooms are keyed by ID; clients by ID; friends and blocks by user ID sets.
- List endpoints format results with numbered lines (`formatNumberedList` in `server/utils.go`).

### Client

- Background goroutine decodes server messages into a channel.
- `handleServerMessage` prints responses and handles prompts (clears partial input when needed).
- `inputLoop` uses readline; slash lines go through `buildCommandMessage`, everything else becomes `SEND_MSG`.
- Terminal actions (`DONE`, `ERR`) and display-only actions (`MESSAGE_ROOM`, `SEND_MSG`, `MESSAGE_FRIEND`) do not block for a follow-up reply.

---

## Limitations

| Area | Status |
|------|--------|
| `SIGN_UP` | Not implemented |
| `SEND_FILE` | Not implemented |
| Persistence | In-memory only; restart clears all state |
| Nickname uniqueness | Not enforced |
| Block list | Stored but not enforced on all code paths |

---

## License

MIT License

## Author

Chima

## Contributing

Issues and pull requests are welcome.
