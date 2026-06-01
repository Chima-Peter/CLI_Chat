# Feature QA-03 — Friend management verification

**Slug:** friend-management-verification  
**Output file:** QA-03-friend-management-verification.md

## Summary

The "Friend management verification" feature allows users to interact with other users on the platform by sending friend requests, accepting or rejecting received requests, cancelling sent requests, managing their friend list, and communicating with friends.

## Scope

- **In scope:**
  - Sending, accepting, rejecting, and cancelling friend requests.
  - Listing friends, sent friend requests, and pending friend requests.
  - Removing friends.
  - Sending messages to friends.
- **Out of scope:**
  - User authentication (covered in `QA-01`).
  - Room management (covered in `QA-02`).
  - Blocking/Unblocking (covered in `QA-05`, although related to friend functionality).

## Functional behaviour

Users can initiate friend requests to other users. Once a request is sent, the recipient can accept or reject it. The sender can cancel a sent request. Users can view their current friends and pending/sent requests. Friends can send private messages to each other, provided they are online and have not blocked each other.

## Technical design

Friend management logic is implemented primarily in the `client` struct and its associated methods in `server/client.go`. The state for friends, pending requests, and sent requests is maintained within the `client` struct (e.g., `friends map[string]struct{}`, `pending_friend_requests map[string]struct{}`, `sent_friend_request map[string]struct{}`). Methods in `server/client.go` handle the state transitions and validation logic (such as checking if users are already friends or blocked).

## Entry points

The entry points are the `client` methods, which are invoked by the server's command handlers (implied as part of the `server` package structure).

## Dependencies

- `server/client.go` (main implementation)
- JSON serialization/deserialization for messages.

## Open questions

None.

## Revision

- Initial draft: scope, behaviour, and technical design from `server/client.go`.

## Evidence index

- `server/client.go:69-113` - Friend request sending logic.
- `server/client.go:115-136` - Friend request acceptance logic.
- `server/client.go:138-165` - Friend request cancellation logic.
- `server/client.go:167-194` - Friend request rejection logic.
- `server/client.go:224-236` - Friend listing logic.
- `server/client.go:238-263` - Friend removal logic.
- `server/client.go:265-321` - Messaging friends logic.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Send friend request](./QA-03_TC-01-send-friend-request.md)
- [Accept friend request](./QA-03_TC-02-accept-friend-request.md)
- [Reject friend request](./QA-03_TC-03-reject-friend-request.md)
- [List friends](./QA-03_TC-04-list-friends.md)
- [Delete friend](./QA-03_TC-05-delete-friend.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



