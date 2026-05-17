package server

import "github.com/chima/CLI_Chat/protocol"

type Message = protocol.Message
type ActionType = protocol.ActionType

const (
	SIGN_UP = protocol.SIGN_UP
	LOGIN   = protocol.LOGIN
	LOGOUT  = protocol.LOGOUT

	CREATE_ROOM                   = protocol.CREATE_ROOM
	SET_ROOM_PASSWORD             = protocol.SET_ROOM_PASSWORD
	JOIN_ROOM                     = protocol.JOIN_ROOM
	LEAVE_ROOM                    = protocol.LEAVE_ROOM
	DELETE_ROOM                   = protocol.DELETE_ROOM
	EDIT_ROOM                     = protocol.EDIT_ROOM
	GET_ROOM_PASSWORD             = protocol.GET_ROOM_PASSWORD
	DELETE_MEMBER                 = protocol.DELETE_MEMBER
	SEND_INVITE_REQUEST           = protocol.SEND_INVITE_REQUEST
	SEE_GROUP_INVITE_REQUEST      = protocol.SEE_GROUP_INVITE_REQUEST
	ACCEPT_GROUP_INVITE_REQUEST   = protocol.ACCEPT_GROUP_INVITE_REQUEST
	DELETE_GROUP_INVITE_REQUEST   = protocol.DELETE_GROUP_INVITE_REQUEST
	GET_ROOM_MEMBERS              = protocol.GET_ROOM_MEMBERS

	LIST_ROOMS    = protocol.LIST_ROOMS
	LIST_MY_ROOMS = protocol.LIST_MY_ROOMS

	SEND_MSG  = protocol.SEND_MSG
	SEND_FILE = protocol.SEND_FILE

	SEND_FRIEND_REQUEST   = protocol.SEND_FRIEND_REQUEST
	ACCEPT_FRIEND_REQUEST = protocol.ACCEPT_FRIEND_REQUEST
	MESSAGE_FRIEND        = protocol.MESSAGE_FRIEND
	GET_FRIENDS           = protocol.GET_FRIENDS
	SEE_FRIEND_REQUEST    = protocol.SEE_FRIEND_REQUEST
	DELETE_FRIEND         = protocol.DELETE_FRIEND
	BLOCK_USER            = protocol.BLOCK_USER
	UNBLOCK_USER          = protocol.UNBLOCK_USER
	GET_USER_STATUS       = protocol.GET_USER_STATUS

	DONE = protocol.DONE
)


// Wiring

// Run() command dispatch — implement CMD_NICK, CMD_JOIN, CMD_ROOMS, CMD_MSG, CMD_QUIT
// JSON action router — replace or complement readInput with a HandleEndpoints-style loop for ActionType values in route_types.go
// Client disconnect — remove client from s.clients, leave rooms, close conn on /quit
// Protocol

// Mixed responses — err / msg use plain ERR: / text lines; send_user_message sends JSON. The client only reads plain lines in read_from_server.go
// readInput vs JSON — pick one protocol for the wire format
// route_types.go actions with no implementation

// Area	Missing
// Auth	SIGN_UP, LOGIN, LOGOUT
// Rooms	EDIT_ROOM, SEE_PENDING_MEMBERS
// Friends	MESSAGE_FRIEND, GET_FRIENDS, DELETE_FRIEND (no methods yet)
// Messaging	SEND_MSG, SEND_FILE (no server handlers)
// All friend/block actions	Methods exist on client / server but nothing calls them
// Logic / behavior

// /join vs CreateRoom — README says join-or-create; JoinRoom only works if the room already exists
// Block list unused — BlockUser does not affect messaging or friend requests yet
// my_rooms not updated — JoinRoom sets cl.room but does not add to my_rooms (used by ListMyRooms / DeleteRoom)
// Missing return after errors — several server methods call cl.err but continue (e.g. SetRoomPassword, JoinRoom, DeleteRoomMember)
// Nickname uniqueness — CMD_NICK handler not implemented; duplicate nicks possible
// No persistence — users, friends, rooms are in-memory only
// Kept on purpose
