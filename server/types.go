package server

import "github.com/chima/CLI_Chat/protocol"

type Message = protocol.Message
type ActionType = protocol.ActionType

const (
	SIGN_UP = protocol.SIGN_UP
	LOGIN   = protocol.LOGIN
	LOGOUT  = protocol.LOGOUT

	CREATE_ROOM                 = protocol.CREATE_ROOM
	SET_ROOM_PASSWORD           = protocol.SET_ROOM_PASSWORD
	JOIN_ROOM                   = protocol.JOIN_ROOM
	LEAVE_ROOM                  = protocol.LEAVE_ROOM
	DELETE_ROOM                 = protocol.DELETE_ROOM
	EDIT_ROOM                   = protocol.EDIT_ROOM
	GET_ROOM_PASSWORD           = protocol.GET_ROOM_PASSWORD
	DELETE_MEMBER               = protocol.DELETE_MEMBER
	SEND_INVITE_REQUEST         = protocol.SEND_INVITE_REQUEST
	SEE_GROUP_INVITE_REQUEST    = protocol.SEE_GROUP_INVITE_REQUEST
	ACCEPT_GROUP_INVITE_REQUEST = protocol.ACCEPT_GROUP_INVITE_REQUEST
	DELETE_GROUP_INVITE_REQUEST = protocol.DELETE_GROUP_INVITE_REQUEST
	GET_ROOM_MEMBERS            = protocol.GET_ROOM_MEMBERS

	LIST_ROOMS           = protocol.LIST_ROOMS
	LIST_MY_ROOMS        = protocol.LIST_MY_ROOMS
	LIST_MY_ROOM_INVITES = protocol.LIST_MY_ROOM_INVITES

	SWITCH_CONTEXT = protocol.SWITCH_CONTEXT
	MESSAGE_ROOM   = protocol.MESSAGE_ROOM
	SEND_MSG       = protocol.SEND_MSG
	SEND_FILE      = protocol.SEND_FILE

	CREATE_FILE_PORT    = protocol.CREATE_FILE_PORT
	FILE_PORT_LISTENING = protocol.FILE_PORT_LISTENING
	UPLOAD_FILE         = protocol.UPLOAD_FILE

	SEND_FRIEND_REQUEST   = protocol.SEND_FRIEND_REQUEST
	ACCEPT_FRIEND_REQUEST = protocol.ACCEPT_FRIEND_REQUEST
	MESSAGE_FRIEND        = protocol.MESSAGE_FRIEND
	GET_FRIENDS           = protocol.GET_FRIENDS
	SEE_FRIEND_REQUEST    = protocol.SEE_FRIEND_REQUEST
	DELETE_FRIEND         = protocol.DELETE_FRIEND
	BLOCK_USER            = protocol.BLOCK_USER
	UNBLOCK_USER          = protocol.UNBLOCK_USER
	GET_USER_STATUS       = protocol.GET_USER_STATUS

	ERR  = protocol.ERR
	DONE = protocol.DONE
)

