package server

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type server struct {
	rooms   map[string]*room
	clients map[string]*client
	mu      sync.RWMutex
}

func InitServer() *server {
	return &server{
		rooms:   make(map[string]*room),
		clients: make(map[string]*client),
	}
}

func (s *server) CreateRoom(cl *client, room_name string) {
	trimmed_name := strings.TrimSpace(room_name)
	if trimmed_name == "" {
		cl.err(fmt.Errorf("Provide a valid room name"))
		return
	}

	if s.isRoomNameTaken(trimmed_name, "") {
		cl.err(fmt.Errorf("%s has already been used. Select another name.", room_name))
		return
	}

	new_room := &room{
		id:         uuid.New().String(),
		name:       trimmed_name,
		is_private: false,
		owner:      cl,
		password:   "",
		members:    make(map[string]*client),
		invites:    make(map[string]*client),
		mu:         sync.RWMutex{},
	}

	new_room.members[cl.id] = cl

	s.mu.Lock()
	s.rooms[new_room.id] = new_room
	s.mu.Unlock()

	cl.my_rooms[new_room.id] = new_room
	cl.room = new_room

	payload := map[string]any{
		"room_id": new_room.id,
		"room":    trimmed_name,
	}

	cl.send_user_message(payload, SET_ROOM_PASSWORD, "Set room password. Leave black for public room: ")
}

func (s *server) SetRoomPassword(cl *client, roomID, roomName, password string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	if room_data.owner.id != cl.id {
		cl.err(fmt.Errorf("You don't have the right to update this room password."))
		return
	}

	room_data.SetRoomPassword(cl, password)
}

func (s *server) EditRoom(cl *client, roomID, roomName, new_name string, max_size *int) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}
	room_data.EditRoom(cl, s, new_name, max_size)
}

func (s *server) JoinRoom(cl *client, roomID, roomName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	room_data.JoinRoom(cl)
}

func (s *server) JoinRoomWithPassword(cl *client, roomID, roomName, password string) {
	trimmed_password := strings.TrimSpace(password)

	if trimmed_password == "" {
		cl.err(fmt.Errorf("Enter password"))
		return
	}

	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	room_data.JoinRoomWithPassword(cl, password)
}

func (s *server) LeaveRoom(cl *client, roomID, roomName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	room_data.LeaveRoom(cl)
}

func (s *server) DeleteRoom(cl *client, roomID, roomName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	room_data.DeleteRoom(cl, s)
}

func (s *server) DeleteRoomMember(cl *client, roomID, roomName, memberID, memberName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	if room_data.owner.id != cl.id {
		cl.err(fmt.Errorf("This action is reserved for only the admin!"))
		return
	}

	member, err := s.resolveUser(memberID, memberName)
	if err != nil {
		cl.err(err)
		return
	}

	_, ok := room_data.members[member.id]
	if !ok {
		cl.err(fmt.Errorf("%s is not a member of room: %s", member.nick, room_data.name))
		return
	}

	room_data.DeleteMember(member)

	member.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Admin deleted you from room: %s", room_data.name))

	room_data.Broadcast(cl, fmt.Sprintf("%s has been deleted from this room", member.nick))

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s has been deleted from this room.", member.nick))
}

func (s *server) GetRoomMembers(cl *client, roomID, roomName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	_, exists := room_data.members[cl.id]
	if !exists {
		cl.err(fmt.Errorf("Only members of a room can see it's members"))
		return
	}

	members := room_data.FetchRoomMembers()

	cl.send_user_message(map[string]any{
		"room_id": room_data.id,
		"room":    room_data.name,
	}, DONE, fmt.Sprintf("Members of this room are: %s", strings.Join(members, ", ")))
}

func (s *server) SendRoomInvite(cl *client, roomID, roomName, memberID, memberName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	if room_data.owner.id != cl.id {
		cl.err(fmt.Errorf("This action is reserved for only the admin!"))
		return
	}

	new_member, err := s.resolveUser(memberID, memberName)
	if err != nil {
		cl.err(err)
		return
	}

	_, exists := room_data.members[new_member.id]
	if exists {
		cl.err(fmt.Errorf("User already part of this room"))
		return
	}

	room_data.SendRoomInvite(new_member, cl)
}

func (s *server) SeePendingRoomInvites(cl *client, roomID, roomName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		cl.err(err)
		return
	}

	if room_data.owner.id != cl.id {
		cl.err(fmt.Errorf("This action is reserved for only the admin!"))
		return
	}

	room_data.SeePendingRoomInvites(cl)
}

func (s *server) AcceptRoomInvite(invitee *client, roomID, roomName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		invitee.err(err)
		return
	}

	room_data.AcceptRoomInvite(invitee)
}

func (s *server) DeclineRoomInvite(invitee *client, roomID, roomName string) {
	room_data, err := s.resolveRoom(roomID, roomName)
	if err != nil {
		invitee.err(err)
		return
	}

	room_data.DeclineRoomInvite(invitee)
}

func (s *server) ListPublicRooms(cl *client) {
	var rooms []map[string]string

	s.mu.RLock()
	for _, room_data := range s.rooms {
		if !room_data.is_private {
			rooms = append(rooms, map[string]string{
				"room_id": room_data.id,
				"room":    room_data.name,
			})
		}
	}
	s.mu.RUnlock()

	if len(rooms) == 0 {
		cl.send_user_message(map[string]any{}, DONE, "No publicly available rooms")
		return
	}

	names := make([]string, 0, len(rooms))
	for _, r := range rooms {
		names = append(names, r["room"])
	}

	cl.send_user_message(map[string]any{"rooms": rooms}, DONE, fmt.Sprintf("Publicly available rooms are: %s", strings.Join(names, ", ")))
}

func (s *server) SendRoomMessage(cl *client, message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		cl.err(fmt.Errorf("message cannot be empty"))
		return
	}
	if cl.room == nil {
		cl.err(fmt.Errorf("join a room before sending a message"))
		return
	}
	cl.room.Broadcast(cl, message)
	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("You: %s", message))
}

func (s *server) ListMyRooms(cl *client) {
	rooms := make([]map[string]string, 0, len(cl.my_rooms))
	names := make([]string, 0, len(cl.my_rooms))

	for _, room_data := range cl.my_rooms {
		rooms = append(rooms, map[string]string{
			"room_id": room_data.id,
			"room":    room_data.name,
		})
		names = append(names, room_data.name)
	}

	cl.send_user_message(map[string]any{"rooms": rooms}, DONE, fmt.Sprintf("You are a member of: %s", strings.Join(names, ", ")))
}

func (s *server) SendFriendRequest(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.SendFriendRequest(target)
}

func (s *server) AcceptFriendRequest(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.AcceptFriendRequest(target)
}

func (s *server) RejectFriendRequest(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.RejectFriendRequest(target)
}

func (s *server) CancelFriendRequest(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.DeleteFriendRequest(target)
}

func (s *server) SeePendingFriendRequests(cl *client) {
	cl.mu.RLock()
	ids := copyIDSet(cl.pending_friend_requests)
	cl.mu.RUnlock()
	cl.FetchPendingRequests(s.usersFromIDs(ids))
}

func (s *server) SeeSentFriendRequests(cl *client) {
	cl.mu.RLock()
	ids := copyIDSet(cl.sent_friend_request)
	cl.mu.RUnlock()
	cl.FetchSentRequests(s.usersFromIDs(ids))
}

func (s *server) GetFriends(cl *client) {
	cl.mu.RLock()
	ids := copyIDSet(cl.friends)
	cl.mu.RUnlock()
	cl.GetFriends(s.usersFromIDs(ids))
}

func (s *server) DeleteFriend(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.DeleteFriend(target)
}

func (s *server) MessageFriend(cl *client, targetID, targetName, message string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.MessageFriend(target, message)
}

func (s *server) BlockUser(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.BlockUser(target)
}

func (s *server) UnblockUser(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.UnblockUser(target)
}

func (s *server) GetUserStatus(cl *client, targetID, targetName string) {
	target, err := s.resolveUser(targetID, targetName)
	if err != nil {
		cl.err(err)
		return
	}
	cl.GetUserStatus(target)
}

func (s *server) LogUserOut(cl *client) {
	cl.setOnline(false)
	for _, room := range cl.my_rooms {
		room.Broadcast(cl, fmt.Sprintf("%s has left the chat.", cl.nick))
		room.mu.Lock()
		delete(room.members, cl.id)
		room.mu.Unlock()
	}
	s.mu.Lock()
	delete(s.clients, cl.id)
	s.mu.Unlock()
	cl.conn.Close()

	log.Println("Client has disconnected:", cl.conn.RemoteAddr().String())
}
