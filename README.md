# CLI Chat

A command-line chat application built with Go that enables real-time communication through a TCP-based server-client architecture.

## Features

- **Multi-client support**: Server handles multiple concurrent client connections
- **Chat rooms**: Create and join chat rooms for organized conversations
- **User nicknames**: Set custom nicknames for identification
- **Room management**: List available rooms and view room members
- **Real-time messaging**: Broadcast messages to all users in a room
- **Command-based interface**: Simple commands for all chat operations

## Project Structure

```
cli_chat/
├── cmd/
│   ├── client/          # Client entry point
│   └── server/          # Server entry point
├── client/              # Client implementation
│   ├── client.go        # Main client connection logic
│   ├── read_from_server.go  # Handle incoming messages
│   └── write_to_server.go   # Handle outgoing messages
├── server/              # Server implementation
│   ├── server.go        # Main server logic
│   ├── client.go        # Client connection handling
│   ├── commands.go      # Command definitions
│   └── rooms.go         # Room and broadcast logic
└── go.mod               # Go module definition
```

## Installation

### Prerequisites
- Go 1.26.3 or higher

### Build from Source

```bash
# Clone the repository
git clone https://github.com/chima/CLI_Chat.git
cd CLI_Chat

# Build the server
go build ./cmd/server

# Build the client
go build ./cmd/client
```

## Usage

### Starting the Server

```bash
./server
```

The server will start listening on `localhost:8888`.

### Connecting a Client

In a separate terminal:

```bash
./client
```

### Available Commands

Once connected to the server, you can use the following commands:

- **/nick `<nickname>`** - Set your nickname
- **/join `<room>`** - Join or create a chat room
- **/rooms** - List all available chat rooms
- **/msg `<message>`** - Send a message to the current room
- **/quit** - Disconnect from the server

### Example Session

```
./client
> /nick alice
> /join general
> /msg Hello everyone!
> /rooms
> /quit
```

## Architecture

### Server
- Listens on port 8888 for TCP connections
- Manages multiple client connections concurrently using goroutines
- Maintains a map of chat rooms
- Processes commands asynchronously through a channel-based command queue
- Broadcasts messages to all members in a room (excluding the sender)

### Client
- Connects to the server on startup
- Handles user input in a separate goroutine
- Receives and displays server messages in a separate goroutine
- Communicates with the server using TCP

## Design Patterns

- **Goroutines**: Concurrent connection handling for scalability
- **Channels**: Command queue for asynchronous message processing
- **Room Broadcasting**: Selective message distribution to room members
- **Anonymous default**: Users start as "anonymous" until they set a nickname

## Future Improvements

- Private direct messages between users
- Room password protection
- Message encryption for secure communication
- Message history stored locally for true anonymity
- Room teardown when empty
- File sharing capabilities
- User authentication
- Configurable server port
- User list in rooms
- Admin commands
- Signal for typing indicators
- Improved error handling and user feedback

## License

MIT License

## Author

Chima

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
