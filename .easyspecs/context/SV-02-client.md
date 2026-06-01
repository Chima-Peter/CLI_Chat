# Service SV-02 — Client

**Slug:** client · **File:** SV-02-client.md

## Summary
The Client service provides the user interface for interacting with the CLI Chat application. It handles connection establishment, session management, and user input processing via a CLI interface.

## Responsibilities
- Establish and maintain a TLS-encrypted connection with the Server service (`client/client.go:11-41`).
- Provide an interactive command-line interface for the user using `readline` (`client/client.go:30-36`).
- Handle user input and map it to protocol messages (`client/session.go:193-256`).
- Manage session state, including pending prompts and terminal display (`client/session.go:22-30`).
- Process asynchronous messages from the server, including file transfer coordination and UI updates (`client/session.go:67-119`).
- Act as a file server/client for file transfers between peers (`client/session.go:277-290`, `client/session.go:292-308`).

## Consumers
The Client service is consumed directly by the end-user through the terminal application built with `github.com/chzyer/readline`. It interfaces with the `protocol` package for message serialization/deserialization.

## Public surface
The primary public interface is the `Connect` method, which initializes the client session.
- `Connect()`: `client/client.go:11`

## Dependencies
- `github.com/chzyer/readline`: CLI interaction.
- `crypto/tls`: Secure connection.
- `github.com/chima/CLI_Chat/protocol`: Message protocol definitions.
- `net.Conn`: Network communication.

## Error model
Errors are handled by displaying them to the user via the `writeDisplay` method, which updates the terminal output without disrupting the active input prompt. Fatal errors (like failed connection establishment) result in disconnection or panic.

## Operational notes
The service uses `SERVER_URL` environment variable to locate the chat server.

## Revision
- Initial draft: Service and method description for Client.

## Evidence index
- `client/client.go:11-41`: Client connection initialization.
- `client/session.go:22-30`: Session structure definition.
- `client/session.go:67-119`: Processing server messages.
- `client/session.go:193-256`: Input loop and user command processing.
- `client/session.go:277-290`: File server initiation.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Connect](./SV-02_ME-05-connect.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



