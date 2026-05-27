package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func (s *server) HandleConn(conn net.Conn) {
	defer conn.Close()

	log.Println("New client is connected:", conn.RemoteAddr().String())

	cl := &client{
		id:                      uuid.New().String(),
		conn:                    conn,
		nick:                    "",
		authenticated:           false,
		room:                    nil,
		current_context:         contextNone,
		current_friend:          nil,
		online:                  false,
		my_rooms:                make(map[string]*room),
		friends:                 make(map[string]struct{}),
		pending_friend_requests: make(map[string]struct{}),
		sent_friend_request:     make(map[string]struct{}),
		blocked_users:           make(map[string]struct{}),
		room_invites:            make(map[string]*room),
		mu:                      sync.RWMutex{},
	}

	s.mu.Lock()
	s.clients[cl.id] = cl
	s.mu.Unlock()

	defer s.LogUserOut(cl)

	cl.setOnline(true)
	cl.promptLogin("Enter your nickname: ")

	decoder := json.NewDecoder(conn)
	for {
		var request Message
		if err := decoder.Decode(&request); err != nil {
			return
		}
		if s.dispatchMessage(cl, &request) {
			return
		}
	}
}

func (s *server) dispatchMessage(cl *client, req *Message) (closeConn bool) {
	p := decodePayload(req.Payload)
	if !cl.authenticated {
		switch req.Action {
		case LOGIN:
			s.handleLogin(cl, p.userName())
			return false
		case LOGOUT:
			return true
		default:
			cl.promptLogin("Log in with a nickname first: ")
			return false
		}
	}

	switch req.Action {
	case SIGN_UP:
		cl.err(fmt.Errorf("sign up is not implemented yet"))
	case LOGIN:
		s.handleLogin(cl, p.userName())
	case LOGOUT:
		cl.send_user_message(map[string]any{}, DONE, "Logged out.")
		return true

	case CREATE_ROOM:
		s.CreateRoom(cl, p.roomName())
	case SET_ROOM_PASSWORD:
		s.SetRoomPassword(cl, p.roomID(), p.roomName(), p.Password)
	case JOIN_ROOM:
		s.JoinRoom(cl, p.roomID(), p.roomName())
	case LEAVE_ROOM:
		s.LeaveRoom(cl, p.roomID(), p.roomName())
	case DELETE_ROOM:
		s.DeleteRoom(cl, p.roomID(), p.roomName())
	case EDIT_ROOM:
		s.EditRoom(cl, p.roomID(), p.roomName(), p.NewRoom, p.MaxSize)
	case GET_ROOM_PASSWORD:
		s.JoinRoomWithPassword(cl, p.roomID(), p.roomName(), p.Password)
	case DELETE_MEMBER:
		s.DeleteRoomMember(cl, p.roomID(), p.roomName(), p.userID(), p.userName())
	case SEND_INVITE_REQUEST:
		s.SendRoomInvite(cl, p.roomID(), p.roomName(), p.userID(), p.userName())
	case SEE_GROUP_INVITE_REQUEST:
		s.SeePendingRoomInvites(cl, p.roomID(), p.roomName())
	case ACCEPT_GROUP_INVITE_REQUEST:
		s.AcceptRoomInvite(cl, p.roomID(), p.roomName())
	case DELETE_GROUP_INVITE_REQUEST:
		s.DeclineRoomInvite(cl, p.roomID(), p.roomName())
	case GET_ROOM_MEMBERS:
		s.GetRoomMembers(cl, p.roomID(), p.roomName())

	case LIST_ROOMS:
		s.ListPublicRooms(cl)
	case LIST_MY_ROOMS:
		s.ListMyRooms(cl)
	case LIST_MY_ROOM_INVITES:
		s.ListMyRoomInvites(cl)

	case SWITCH_CONTEXT:
		switch p.contextType() {
		case contextFriend:
			s.SwitchContext(cl, contextFriend, p.friendName())
		case contextRoom:
			s.SwitchContext(cl, contextRoom, p.roomName())
		default:
			s.SwitchContext(cl, p.contextType(), p.roomName())
		}

	case MESSAGE_ROOM:
		s.SendRoomMessage(cl, p.roomName(), p.Message)
	case SEND_MSG:
		s.SendContextMessage(cl, p.Message)
	case SEND_FILE:
		s.requestFriendFilePort(cl, p.friendName(), p.filePath())
	case FILE_PORT_LISTENING:
		s.handleFilePortListening(cl, req.Payload)
	case SEND_FRIEND_REQUEST:
		s.SendFriendRequest(cl, p.userID(), p.userName())
	case ACCEPT_FRIEND_REQUEST:
		s.AcceptFriendRequest(cl, p.userID(), p.userName())
	case MESSAGE_FRIEND:
		s.MessageFriend(cl, p.userID(), p.userName(), p.Message)
	case GET_FRIENDS:
		s.GetFriends(cl)
	case SEE_FRIEND_REQUEST:
		s.SeePendingFriendRequests(cl)
	case DELETE_FRIEND:
		s.DeleteFriend(cl, p.userID(), p.userName())
	case BLOCK_USER:
		s.BlockUser(cl, p.userID(), p.userName())
	case UNBLOCK_USER:
		s.UnblockUser(cl, p.userID(), p.userName())
	case GET_USER_STATUS:
		s.GetUserStatus(cl, p.userID(), p.userName())

	default:
		cl.err(fmt.Errorf("unknown action: %d", req.Action))
	}

	return false
}

func (cl *client) promptLogin(msg string) {
	cl.send_user_message(map[string]any{}, LOGIN, msg)
}

func (s *server) handleLogin(cl *client, nick string) {
	if cl.authenticated {
		cl.promptLogin("You are already logged in. Log out first: ")
		return
	}
	nick = strings.TrimSpace(nick)
	if nick == "" {
		cl.promptLogin("Nickname cannot be empty. Enter a valid nickname: ")
		return
	}
	if strings.EqualFold(nick, "anonymous") {
		cl.promptLogin("Nickname cannot be anonymous. Enter a valid nickname: ")
		return
	}
	if s.isNickTaken(nick, cl.id) {
		cl.promptLogin("Nickname is already taken. Enter a valid nickname: ")
		return
	}

	cl.mu.Lock()
	cl.nick = nick
	cl.authenticated = true
	cl.mu.Unlock()

	cl.send_user_message(
		map[string]any{"user_id": cl.id, "nick": nick},
		DONE,
		fmt.Sprintf("Logged in as %s.", nick),
	)
	log.Printf("Requesting file server for user %s", nick)
	cl.send_user_message(map[string]any{}, CREATE_FILE_PORT, "")
}
