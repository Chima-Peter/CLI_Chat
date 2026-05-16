package server

import (
	"encoding/json"
	"fmt"
)

func (cl *client) HandleEndpoints() {
	decoder := json.NewDecoder(cl.conn)

	for {
		var request Message
		err := decoder.Decode(&request)
		if err != nil {
			return
		}

		switch request.Action {
		// Authentication endpoints
		case SIGN_UP:
			fmt.Println("Register new user account")
		case LOGIN:
			fmt.Println("Authenticate user with credentials")
		case LOGOUT:
			fmt.Println("Sign out user session")

		// Room management endpoints
		case CREATE_ROOM:
			fmt.Println("Create a new chat room")
		case SET_ROOM_PASSWORD:
			fmt.Println("(Owner only) Set password protection for room")
		case JOIN_ROOM:
			fmt.Println("Join an existing chat room")
		case LEAVE_ROOM:
			fmt.Println("Exit from a room")
		case DELETE_ROOM:
			fmt.Println("(Owner only) Delete room permanently")
		case EDIT_ROOM:
			fmt.Println("(Owner only) Modify room settings and configuration")
		case GET_ROOM_PASSWORD:
			fmt.Println("Retrieve room password for authentication")
		case SEE_PENDING_MEMBERS:
			fmt.Println("(Owner only) List members waiting for approval")
		case DELETE_MEMBER:
			fmt.Println("(Owner only) Remove member from room")
		case GET_ROOM_MEMBERS:
			fmt.Println("Retrieve list of room members")

		// Room invitation endpoints
		case SEND_INVITE_REQUEST:
			fmt.Println("(Owner only) Send invitation request to join room")
		case SEE_GROUP_INVITE_REQUEST:
			fmt.Println("(Owner only) View pending group invitations")
		case ACCEPT_GROUP_INVITE_REQUEST:
			fmt.Println("Accept invitation to join group")
		case DELETE_GROUP_INVITE_REQUEST:
			fmt.Println("Decline or delete invitation")

		// Room listing endpoints
		case LIST_ROOMS:
			fmt.Println("List all publicly available rooms")
		case LIST_MY_ROOMS:
			fmt.Println("List rooms user is member of")

		// Messaging endpoints
		case SEND_MSG:
			fmt.Println("Send message to room or user")
		case SEND_FILE:
			fmt.Println("Upload and send file to room or user")

		// Friend management endpoints
		case SEND_FRIEND_REQUEST:
			fmt.Println("Send friend request to user")
		case ACCEPT_FRIEND_REQUEST:
			fmt.Println("Accept incoming friend request")
		case MESSAGE_FRIEND:
			fmt.Println("Send direct message to friend")
		case GET_FRIENDS:
			fmt.Println("Retrieve list of friends")
		case SEE_FRIEND_REQUEST:
			fmt.Println("View pending friend requests")
		case DELETE_FRIEND:
			fmt.Println("Remove user from friend list")

		// User interaction endpoints
		case BLOCK_USER:
			fmt.Println("Block user from messaging")
		case UNBLOCK_USER:
			fmt.Println("Unblock user to allow messaging")
		case GET_USER_STATUS:
			fmt.Println("Check if user is online")

		default:
			fmt.Println("Unknown action type")
		}
	}
}
