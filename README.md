# CLI Chat

A terminal-based chat platform written in Go: a **TLS** chat server and an interactive client that share a **line-delimited JSON protocol**. User commands are grouped like REST-style routes (`/room/join`, `/friend/add`, `/chat/dm`). The server maintains an **active chat context** so plain text is routed to the correct room or friend.

This document is intended for **operators, contributors, and end users** who want to run, extend, or deploy the platform publicly.

---

## Table of contents

- [Overview](#overview)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Client commands](#client-commands)
- [Authentication](#authentication)
- [Messaging](#messaging)
- [Friend-to-friend file transfer](#friend-to-friend-file-transfer)
- [Wire protocol](#wire-protocol)
- [Architecture](#architecture)
- [Security and public deployment](#security-and-public-deployment)
- [Known limitations](#known-limitations)
- [Project layout](#project-layout)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

CLI Chat is a **multi-user, in-memory** chat system over **TCP + TLS**:

| Layer | Detail |
|-------|--------|
| Transport | TLS 1.x (`crypto/tls`) on a single long-lived connection per client |
| Framing | One JSON object (`protocol.Message`) per line (`encoding/json` decoder) |
| Identity | Nickname chosen at login (no password); UUID per connection |
| State | Rooms, membership, friends, blocks, and invites live **only in server RAM** |

The client uses [readline](https://github.com/chzyer/readline) for interactive input. The server handles **one goroutine per connection** and dispatches by `action` integer.

**Not a production-hardened SaaS today:** there is no account system, persistence, rate limiting, or certificate pinning. See [Security and public deployment](#security-and-public-deployment) before exposing a host to the internet.

---

## Features

- **Concurrent TLS server** — accept loop + `HandleConn` per client
- **JSON wire protocol** — shared `protocol` package; action IDs must match between client and server builds
- **Mandatory login** — nickname prompt on connect; reserved name `anonymous` rejected; uniqueness among online users
- **Rooms** — create, join, leave, delete, rename, max size, optional password, invites, member list, owner kick
- **Friends** — requests, accept, list, remove, direct messages (online only)
- **Context switching** — `/switch room|friend <name>` sets where plain text goes
- **Explicit chat** — `/chat/room` and `/chat/dm` without changing context
- **Server-driven prompts** — login, room password on create/join
- **Numbered list responses** — rooms, members, invites, friends (`1. …`, `2. …`)
- **User blocks** — block/unblock; enforced on friend requests and DMs
- **Peer file transfer** — friends can send files from `~/Downloads` via a secondary TLS listener on the receiver’s machine (see [Friend-to-friend file transfer](#friend-to-friend-file-transfer))

**Stubs / partial:** `SIGN_UP`, `/chat/file` (legacy path field), `FILE_SERVER_CLOSED` (defined, unused).

---

## Requirements

- **Go 1.26.3+** (see `go.mod`)
- A terminal with ANSI support (client clears screen on `clear`)
- **Windows / macOS / Linux** — file save path uses `~/Downloads`
- For remote use: network path from sender → receiver’s **file listener** host/port (often requires LAN or explicit port forwarding)

---

## Installation

```bash
git clone https://github.com/chima/CLI_Chat.git
cd CLI_Chat

go build -o server ./cmd/server
go build -o client ./cmd/client
```

On Windows, binaries are commonly `server.exe` and `client.exe`.

---

## Configuration

| Variable | Used by | Default | Purpose |
|----------|---------|---------|---------|
| `PORT` | server | `8080` | TLS listen port (`:PORT`) |
| `SERVER_URL` | client | `localhost:8080` | Host (and port) for `tls.Dial` |

Examples:

```bash
# Server on port 9000
PORT=9000 ./server

# Client to a remote host
SERVER_URL=chat.example.com:9000 ./client
```

PowerShell:

```powershell
$env:PORT = "9000"
.\server.exe
$env:SERVER_URL = "chat.example.com:9000"
.\client.exe
```

---

## Usage

**Terminal 1 — start the server:**

```bash
./server
```

**Terminal 2 — start a client:**

```bash
./client
```

The client prints a short banner; type **`/help`** for the full command list.

### Quick start (two users)

**Alice:**

```
alice
/room/create general
# (prompt) leave password blank for a public room
/room/join general
/switch room general
Hello everyone!
```

**Bob** (second terminal):

```
bob
/room/join general
/switch room general
Hi Alice!
```

> Joining a room does **not** set chat context. Use `/switch room <name>` before plain text, or `/chat/room <name> <message>` for a one-off message.

### Direct messages

```
/friend/add bob
# bob runs: /friend/accept alice
/switch friend bob
Hey Bob!
```

Or without switching: `/chat/dm bob Hey Bob!`

### Disconnect

`/quit` or Ctrl+C / EOF on the client sends `LOGOUT` and closes the session.

---

## Client commands

Slash commands are built in `client/message.go` and sent as `protocol.Message` values. `/help` is handled only on the client.

### Auth

| Command | Action | Notes |
|---------|--------|-------|
| *(on connect)* | `LOGIN` | Server prompts; only login/plain nickname until `DONE` |
| `/login <name>` | `LOGIN` | Same as typing a nickname at the prompt |
| `/signup <name>` | `SIGN_UP` | **Not implemented** on server |
| `/quit` | `LOGOUT` | Disconnect |

### Room

| Command | Action |
|---------|--------|
| `/room/create <name>` | `CREATE_ROOM` |
| `/room/join <name>` | `JOIN_ROOM` |
| `/room/leave [name]` | `LEAVE_ROOM` |
| `/room/delete <name>` | `DELETE_ROOM` |
| `/room/edit <room> <new_name> [max]` | `EDIT_ROOM` |
| `/room/members <room>` | `GET_ROOM_MEMBERS` |
| `/room/kick <room> <user>` | `DELETE_MEMBER` |
| `/room/invite <room> <user>` | `SEND_INVITE_REQUEST` |
| `/room/invites [room]` | `SEE_GROUP_INVITE_REQUEST` |
| `/room/list` | `LIST_ROOMS` |
| `/room/mine` | `LIST_MY_ROOMS` |

### Invite

| Command | Action |
|---------|--------|
| `/invite/mine` | `LIST_MY_ROOM_INVITES` |
| `/invite/accept <room>` | `ACCEPT_GROUP_INVITE_REQUEST` |
| `/invite/decline <room>` | `DELETE_GROUP_INVITE_REQUEST` |

### Chat

| Command | Action | Notes |
|---------|--------|-------|
| `/switch room <name>` | `SWITCH_CONTEXT` | Active context = room |
| `/switch friend <name>` | `SWITCH_CONTEXT` | Active context = friend |
| *plain text* | `SEND_MSG` | Requires active context |
| `/chat/room <room> <message>` | `MESSAGE_ROOM` | Does not change context |
| `/chat/dm <friend> <message>` | `MESSAGE_FRIEND` | Does not change context |
| `/chat/file <path>` | `SEND_FILE` | Legacy; prefer `/file` (see below) |

### Friend

| Command | Action |
|---------|--------|
| `/friend/add <user>` | `SEND_FRIEND_REQUEST` |
| `/friend/accept <user>` | `ACCEPT_FRIEND_REQUEST` |
| `/friend/remove <user>` | `DELETE_FRIEND` |
| `/friend/requests` | `SEE_FRIEND_REQUEST` |
| `/friend/list` | `GET_FRIENDS` |

### User

| Command | Action |
|---------|--------|
| `/user/block <user>` | `BLOCK_USER` |
| `/user/unblock <user>` | `UNBLOCK_USER` |
| `/user/status <user>` | `GET_USER_STATUS` |

### Files

| Command | Action | Notes |
|---------|--------|-------|
| `/file <friend> <filename>` | `SEND_FILE` | File must exist in **`~/Downloads`**; zipped client-side |

Names (rooms, users, friends) **cannot contain spaces**.

---

## Authentication

There is **no password or token**. A client proves identity only by picking a unique nickname after connect.

```
Client connects
  → Server: LOGIN "Enter your nickname:"
  → Client: prompt mode (slash commands disabled until login completes)

User sends nickname (or /login <name>)
  → Server: handleLogin
       • empty / "anonymous" (case-insensitive) / taken → LOGIN + error, stay in prompt
       • valid → DONE "Logged in as …"
       → Server: CREATE_FILE_PORT (receiver starts file listener)

Any other action before login
  → Server: LOGIN "Log in with a nickname first: "
```

While unauthenticated, the client treats server messages with action `LOGIN` as prompts and routes the next line to `LOGIN` with merged payload.

**Login rules:**

- Nickname non-empty
- Not `anonymous` (any casing)
- Unique among **authenticated** clients (by nickname string)

---

## Messaging

Three paths:

```
Plain text
  → SEND_MSG { "message": "..." }
  → SendContextMessage
       • context "room"  → broadcast to active room
       • context "friend" → MessageFriend (DM)
       • no context → error

/chat/room <room> <message>
  → MESSAGE_ROOM
  → SendRoomMessage (membership required)

/chat/dm <friend> <message>
  → MESSAGE_FRIEND
  → MessageFriend (must be friends; friend online; block checks)
```

### Server context (per connection)

| Field | Meaning |
|-------|---------|
| `current_context` | `""`, `"room"`, or `"friend"` |
| `room` | Active room when context is `room` |
| `current_friend` | Active friend when context is `friend` |

`my_rooms` (membership) is separate: you can belong to many rooms but only one **active** plain-text target.

### Incoming traffic

- **Room:** other members receive `MESSAGE_ROOM` with `response_msg` for display.
- **DM:** recipient receives `MESSAGE_FRIEND`.

---

## Friend-to-friend file transfer

Files are sent **friend-to-friend**, not in room chat. The chat server **coordinates** endpoints; payload bytes go **directly** from sender to receiver over a **second TLS connection**.

### Flow

```
1. Receiver logs in
   → Server sends CREATE_FILE_PORT
   → Client starts TLS listener on ephemeral port (:0)
   → Client sends FILE_PORT_LISTENING { host, port } to server

2. Sender: /file <friend> <filename>
   → Client zips ~/Downloads/<filename> → ~/Downloads/<base>.zip
   → SEND_FILE { filepath, username } to server
   → Server checks friendship; if receiver ready, sends sender FILE_PORT_LISTENING
      { host, port, nick, friend, filepath }
   → Sender dials receiver:port with TLS, sends UPLOAD_FILE metadata + raw zip bytes

3. Receiver file server saves to ~/Downloads/<name>
```

If the receiver’s file port is not ready yet, the server **queues** the request and delivers `FILE_PORT_LISTENING` when the receiver registers.

### File rules (client)

- Source file must be in **`~/Downloads`** (basename only; path traversal stripped)
- Max source size **100 MiB** before zipping
- Blocked extensions: `.exe`, `.bat`, `.ps1`, `.dll`, etc. (see `client/file-parsing.go`)
- Many already-compressed types (`.zip`, `.jpg`, `.mp4`, …) are rejected for zipping
- Receiver writes to **`~/Downloads`** with `filepath.Base` on the transferred name

### Network note

The file listener advertises **`127.0.0.1`** when bound to all interfaces (`listenerEndpoint`). **Remote senders cannot reach another machine’s localhost.** Public or LAN deployment requires advertising a reachable host (today this is development-oriented; see limitations).

---

## Wire protocol

### Message shape

One JSON object per line (newline-delimited JSON over TLS):

```json
{
  "action": 1,
  "response_msg": "Enter your nickname: ",
  "payload": {}
}
```

| Field | Type | Description |
|-------|------|-------------|
| `action` | int | `ActionType` enum (see table below) |
| `response_msg` | string | Human-readable status or chat line |
| `payload` | object | Action-specific fields (may be `{}`) |

### Action IDs (`protocol/message.go`)

Client and server **must be built from the same protocol revision**. Values are `iota` order:

| ID | Constant | Category |
|----|----------|----------|
| 0 | `SIGN_UP` | Auth (stub) |
| 1 | `LOGIN` | Auth |
| 2 | `LOGOUT` | Auth |
| 3 | `CREATE_ROOM` | Room |
| 4 | `SET_ROOM_PASSWORD` | Room prompt |
| 5 | `JOIN_ROOM` | Room |
| 6 | `LEAVE_ROOM` | Room |
| 7 | `DELETE_ROOM` | Room |
| 8 | `EDIT_ROOM` | Room |
| 9 | `GET_ROOM_PASSWORD` | Room prompt |
| 10 | `DELETE_MEMBER` | Room |
| 11 | `SEND_INVITE_REQUEST` | Room |
| 12 | `SEE_GROUP_INVITE_REQUEST` | Room |
| 13 | `ACCEPT_GROUP_INVITE_REQUEST` | Room |
| 14 | `DELETE_GROUP_INVITE_REQUEST` | Room |
| 15 | `GET_ROOM_MEMBERS` | Room |
| 16 | `LIST_ROOMS` | Room |
| 17 | `LIST_MY_ROOMS` | Room |
| 18 | `LIST_MY_ROOM_INVITES` | Room |
| 19 | `SWITCH_CONTEXT` | Chat |
| 20 | `MESSAGE_ROOM` | Chat |
| 21 | `SEND_MSG` | Chat |
| 22 | `SEND_FILE` | Chat / files |
| 23 | `SEND_FRIEND_REQUEST` | Friend |
| 24 | `ACCEPT_FRIEND_REQUEST` | Friend |
| 25 | `MESSAGE_FRIEND` | Friend |
| 26 | `GET_FRIENDS` | Friend |
| 27 | `SEE_FRIEND_REQUEST` | Friend |
| 28 | `DELETE_FRIEND` | Friend |
| 29 | `BLOCK_USER` | User |
| 30 | `UNBLOCK_USER` | User |
| 31 | `GET_USER_STATUS` | User |
| 32 | `CREATE_FILE_PORT` | File |
| 33 | `FILE_PORT_LISTENING` | File |
| 34 | `UPLOAD_FILE` | File (secondary connection) |
| 35 | `FILE_SERVER_CLOSED` | File (unused) |
| 36 | `ERR` | Response |
| 37 | `DONE` | Response |

### Common payload fields

| JSON field | Usage |
|------------|--------|
| `username`, `nick` | Login, user targets |
| `user_id` | Returned on successful login |
| `room`, `room_id`, `room_name` | Room operations |
| `new_room`, `max_size` | Edit room |
| `password` | Private room create/join |
| `message` | Chat body |
| `context` | `SWITCH_CONTEXT`: `"room"` or `"friend"` |
| `friend_name` | Switch to friend |
| `filepath`, `username` | `SEND_FILE` |
| `host`, `port` | File listener registration / delivery |
| `from`, `from_user_id` | Incoming chat metadata |

### Server-driven prompts

If the server reply action is **not** `DONE`, `ERR`, `MESSAGE_ROOM`, `SEND_MSG`, or `MESSAGE_FRIEND`, the client enters **prompt mode**: the next input line is sent with the **same action** and merged payload. Slash commands are not parsed until the prompt completes (`DONE` / `ERR`).

| Prompt action | When |
|---------------|------|
| `LOGIN` | Connect, failed nickname, or action before login |
| `SET_ROOM_PASSWORD` | After `/room/create` |
| `GET_ROOM_PASSWORD` | Joining a private room |

### Client terminal vs display actions

- **Terminal:** `DONE`, `ERR` — end prompt / show message
- **Display-only:** `MESSAGE_ROOM`, `SEND_MSG`, `MESSAGE_FRIEND` — print `response_msg`, no reply required
- **Special:** `CREATE_FILE_PORT`, `FILE_PORT_LISTENING` — handled in `client/session.go` without blocking the input loop

---

## Architecture

```
┌─────────────┐     TLS (PORT)      ┌──────────────────────────────────┐
│   client    │◄───────────────────►│  server                          │
│  readline   │  JSON lines         │  InitServer()                    │
│  session    │                     │  HandleConn → dispatchMessage    │
└──────┬──────┘                     │  rooms, clients, pending files   │
       │                            └──────────────────────────────────┘
       │ TLS (ephemeral port)
       ▼
┌─────────────┐   FILE_PORT_LISTENING coordinates host:port
│ file client │◄────────────────────────────────────────────── sender
│ file server │   UPLOAD_FILE + raw bytes on secondary conn
└─────────────┘
```

### Server (`server/`)

- **`HandleConn`** — register `client`, `promptLogin`, read loop, `dispatchMessage`
- **Pre-auth gate** — only `LOGIN` and `LOGOUT` until `authenticated`
- **Handlers** — `server.go` (rooms, context, friends, files), `rooms.go` (membership, broadcast, invites), `client.go` (per-user state, blocks, friend ops)
- **Replies** — `send_user_message` → JSON encoder on the TLS connection
- **IDs** — rooms and clients use UUID strings; nicknames are display names

### Client (`client/`)

- **`Connect`** — TLS dial, readline, `runSession`
- **`readFromServer`** — decode loop → channel
- **`handleServerMessage`** — prompts, display, file port side effects
- **`inputLoop`** — prompts, `/file`, slash commands, plain `SEND_MSG`
- **`tls_config`** — shared dev certificate generator for server and file listeners

### Room model (`server/rooms.go`)

| Field | Purpose |
|-------|---------|
| `id`, `name` | Identity |
| `owner` | Creator; can kick, see outbound invites |
| `is_private`, `password` | Gated join |
| `max_size` | `0` = unlimited |
| `members`, `invites` | Membership and pending invites |

Public rooms appear in `LIST_ROOMS`. Private rooms require password on join or owner invite flow.

---

## Security and public deployment

Read this section before exposing CLI Chat on the internet.

### Current security posture

| Topic | Behavior |
|-------|----------|
| TLS certificates | **Generated at runtime** in `tls_config.TLSDevConfig()` (self-signed, localhost SANs, 1-year) |
| Client TLS verify | **`InsecureSkipVerify: true`** — no certificate validation |
| Authentication | **Nickname only** — anyone who guesses a name can impersonate after disconnect |
| Authorization | Room owner rules + friendship for DMs/files; **no global admin** |
| Persistence | None — restart wipes users, rooms, friends |
| Rate limiting | None |
| Input validation | Basic string checks; not hardened against abuse |
| File transfer | Friend-only; extension blocklist; size cap; saves to Downloads |

### Recommendations for public use

1. **Replace dev TLS** — Load real certificates (Let’s Encrypt, ACM, etc.) and remove `InsecureSkipVerify` on clients; pin or trust a proper CA.
2. **Add real accounts** — Passwords, OAuth, or mTLS before advertising as a public service.
3. **Persist or replicate state** — If you need reliability across restarts.
4. **Firewall** — Expose only `PORT`; file ports are dynamic and peer-to-peer.
5. **Fix file reachability** — Publish a routable `host` for `FILE_PORT_LISTENING` when clients are not on the same machine.
6. **Operational limits** — Connection caps, message size limits, rate limits, logging, and monitoring.
7. **Threat model** — Treat chat and file content as **untrusted user input**; do not run the server as root.

### Privacy

Messages and metadata are visible to the server process in memory. DMs are not end-to-end encrypted. File bytes pass through the receiver’s listener, not the central server payload stream (after coordination).

---

## Known limitations

| Area | Status |
|------|--------|
| `SIGN_UP` | Not implemented |
| `/chat/file` | Sends `path` field; server expects `filepath` + friend — use `/file` |
| `FILE_SERVER_CLOSED` | Defined, not used |
| Persistence | In-memory only |
| Remote file send | Receiver often advertises `127.0.0.1` |
| Room messages | No history; offline users miss traffic |
| DMs | Online friends only |
| Nicknames | Reusable after disconnect; no account recovery |
| Server bootstrap banner | Outdated sample commands in `cmd/server/main.go` |
| Horizontal scale | Single process; no shared state between instances |

---

## Project layout

```
CLI_Chat/
├── cmd/
│   ├── client/main.go       # Client entry (Connect)
│   └── server/main.go       # TLS listener, PORT env
├── protocol/
│   └── message.go           # ActionType + Message (shared contract)
├── tls_config/
│   └── config.go            # Dev TLS certificate generation
├── client/
│   ├── client.go            # Dial and session bootstrap
│   ├── session.go           # Readline, prompts, file port hooks
│   ├── message.go           # Slash commands → protocol
│   ├── help.go              # /help and banner
│   ├── file-parsing.go      # Zip + extension policy
│   ├── file-server.go       # Incoming file TLS server
│   └── file-client.go       # Outgoing file upload
├── server/
│   ├── server.go            # Room, context, friend, file coordination
│   ├── handle_conn.go       # Per-connection loop and dispatch
│   ├── client.go            # Per-client state and social ops
│   ├── rooms.go             # Room struct, broadcast, invites
│   ├── utils.go             # Payload parsing, list formatting
│   └── types.go             # Re-exports protocol constants
├── go.mod
├── LICENSE
└── README.md
```

---

## Contributing

Issues and pull requests are welcome. When changing `protocol/message.go`, update this README’s action table and note the breaking change in the PR.

Suggested checks before submitting:

```bash
go build ./...
go vet ./...
```

---

## License

[MIT License](LICENSE)

## Author

Chima
