# Method ME-31 — GetFriends

**Service:** SV-01 · **File:** SV-01_ME-31-get-friends.md

## Summary
The `GetFriends` method allows a client to retrieve a list of their current friends.

## Operation
The `GetFriends` function is defined in the `server` struct and is invoked to fetch the friend IDs associated with the calling client, look up the user details for those IDs, and send the formatted list back to the client.

## Request / inputs
The method accepts a `cl *client` parameter, representing the client requesting their friends list.

## Response / outputs
The method does not return a value directly. Instead, it interacts with the client's connection via `cl.GetFriends` to send a JSON message containing the list of friends.

## Auth and permissions
The method operates in the context of an authenticated client connection. No specific explicit authorization check is performed within the `GetFriends` method itself, as the client object `cl` is assumed to be authorized when `GetFriends` is called on the server.

## Idempotency and concurrency
The method uses `cl.mu.RLock()` and `cl.mu.RUnlock()` to safely copy the `cl.friends` map, ensuring concurrency safety during the retrieval of friend IDs. The operation is idempotent as it only retrieves and displays existing information without modifying the state.

## Errors
No specific error cases are explicitly handled within `GetFriends` beyond standard server operations.

## Implementation notes
The method relies on `s.usersFromIDs(ids)` to resolve the user details for the retrieved friend IDs and then invokes `cl.GetFriends(users)` to send the data back to the client.

## Revision
- Initial draft: contract and handler from implementation.

## Evidence index
- `server/server.go:626-631`: The `GetFriends` method definition in the `server` struct.
- `server/client.go:224-236`: The `GetFriends` method definition in the `client` struct, which sends the response to the client.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [Server](./SV-01-server.md)
- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



