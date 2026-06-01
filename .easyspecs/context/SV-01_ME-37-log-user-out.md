# Method ME-37 — LogUserOut

**Service:** SV-01 · **File:** SV-01_ME-37-log-user-out.md

## Summary

The `LogUserOut` method handles the process of a user logging out of the server. It performs necessary cleanup, including updating the user's online status, notifying rooms, removing the user from the server's active clients list, and closing the connection.

## Operation

`func (s *server) LogUserOut(cl *client)`

This method is implemented in `server/server.go`.

## Request / inputs

- **`cl`**: A pointer to the `client` object representing the user to be logged out.

## Response / outputs

This method does not return any values; its purpose is to perform side-effect operations related to user disconnection.

## Auth and permissions

There are no explicit authentication checks performed within this method, as it is assumed that the `client` object passed is already connected and authenticated.

## Idempotency and concurrency

- **Idempotency**: This method involves cleanup actions such as deleting entries from maps and closing connections. While it is not strictly idempotent in a functional sense (e.g., closing a connection that is already closed might have implications depending on the `net.Conn` implementation), the use of map deletions and `setOnline(false)` generally makes subsequent calls safe from causing inconsistent states.
- **Concurrency**: The method utilizes `room.mu.Lock()` when modifying room members and `s.mu.Lock()` when modifying the server's clients map, ensuring thread-safe operations during these critical sections.

## Errors

There are no explicit error returns. Errors related to connection closing or state updates are not bubbled up, and the method relies on logging for visibility.

## Implementation notes

- Sets the client's online status to `false`.
- Iterates over the client's current rooms, broadcasts a departure message to each room, and removes the client from the room's members list.
- Locks the server's mutex to remove the client from the `s.clients` map.
- Closes the client's network connection (`cl.conn.Close()`).
- Logs the disconnection event and the closing of the file server (if applicable).

## Revision

- Initial draft: contract and handler from implementation.

## Evidence index

- `server/server.go:678-693`: Implementation of `LogUserOut` method.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



