package server

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type pendingFileSend struct {
	senderID string
	filePath string
}

type server struct {
	rooms               map[string]*room
	clients             map[string]*client
	pendingFileRequests map[string][]pendingFileSend // receiver id -> pending senders
	mu                  sync.RWMutex
}

func InitServer() *server {
	return &server{
		rooms:               make(map[string]*room),
		clients:             make(map[string]*client),
		pendingFileRequests: make(map[string][]pendingFileSend),
	}
}

// ROOM METHODS

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

	cl.send_user_message(payload, SET_ROOM_PASSWORD, "Set room password. Leave blank for public room: ")
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
		cl.err(fmt.Errorf("Provide a valid password"))
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

	member.send_user_message(map[string]any{}, DONE, fmt.Sprintf("[%s] Admin removed you from room: %s", room_data.name, room_data.name))

	room_data.Broadcast(cl, fmt.Sprintf("%s has been removed from this room", member.nick))

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("You have removed %s from this room.", member.nick))
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
	}, DONE, formatNumberedList("these are the room members:", members))
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
		cl.err(fmt.Errorf("%s is already part of this room", new_member.nick))
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

	cl.send_user_message(map[string]any{"rooms": rooms}, DONE, formatNumberedList("these are the publicly available rooms:", names))
}

func (cl *client) roomByName(name string) (*room, bool) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	for _, room_data := range cl.my_rooms {
		if room_data.name == name {
			return room_data, true
		}
	}
	return nil, false
}

func (s *server) SwitchContext(cl *client, contextType, name string) {
	contextType = strings.TrimSpace(strings.ToLower(contextType))
	name = strings.TrimSpace(name)
	if name == "" {
		cl.err(fmt.Errorf("usage: /switch <room|friend> <name>"))
		return
	}

	switch contextType {
	case contextRoom:
		room_data, ok := cl.roomByName(name)
		if !ok {
			room_data, err := s.resolveRoom("", name)
			if err != nil {
				cl.err(fmt.Errorf("you are not in room: %s", name))
				return
			}
			cl.mu.RLock()
			_, ok = cl.my_rooms[room_data.id]
			cl.mu.RUnlock()
			if !ok {
				cl.err(fmt.Errorf("you are not in room: %s", name))
				return
			}
		}

		cl.mu.Lock()
		cl.current_context = contextRoom
		cl.room = room_data
		cl.current_friend = nil
		cl.mu.Unlock()

		cl.send_user_message(map[string]any{
			"context": contextRoom,
			"room":    room_data.name,
		}, DONE, fmt.Sprintf("Switched to room: %s", room_data.name))

	case contextFriend:
		friend, err := s.resolveUser("", name)
		if err != nil {
			cl.err(err)
			return
		}

		cl.mu.RLock()
		isFriend := hasID(cl.friends, friend.id)
		cl.mu.RUnlock()
		if !isFriend {
			cl.err(fmt.Errorf("you are not friends with %s", friend.nick))
			return
		}

		cl.mu.Lock()
		cl.current_context = contextFriend
		cl.current_friend = friend
		cl.room = nil
		cl.mu.Unlock()

		cl.send_user_message(map[string]any{
			"context": contextFriend,
			"nick":    friend.nick,
		}, DONE, fmt.Sprintf("Switched to friend: %s", friend.nick))

	default:
		cl.err(fmt.Errorf("context must be room or friend"))
	}
}

func (s *server) SendContextMessage(cl *client, message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		cl.err(fmt.Errorf("message cannot be empty"))
		return
	}

	cl.mu.RLock()
	ctx := cl.current_context
	room_data := cl.room
	friend := cl.current_friend
	cl.mu.RUnlock()

	switch ctx {
	case contextRoom:
		if room_data == nil {
			cl.err(fmt.Errorf("use /switch room <name> before sending a message"))
			return
		}
		room_data.Broadcast(cl, message)
		cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("You: %s", message))

	case contextFriend:
		if friend == nil {
			cl.err(fmt.Errorf("use /switch friend <name> before sending a message"))
			return
		}
		cl.MessageFriend(friend, message)

	default:
		cl.err(fmt.Errorf("use /switch room|friend <name> before sending a message"))
	}
}

func (s *server) requestFriendFilePort(cl *client, friendName, filePath string) {
	friendName = strings.TrimSpace(friendName)
	if friendName == "" {
		cl.err(fmt.Errorf("friend name is required"))
		return
	}
	if strings.TrimSpace(filePath) == "" {
		cl.err(fmt.Errorf("file path is required"))
		return
	}

	friend, err := s.resolveUser("", friendName)
	if err != nil {
		cl.err(err)
		return
	}

	friend.mu.RLock()
	requesterIsFriend := hasID(friend.friends, cl.id)
	friend.mu.RUnlock()
	if !requesterIsFriend {
		cl.err(fmt.Errorf("you are not in %s's friend list", friend.nick))
		return
	}

	host, port := friend.GetFriendFilePort()
	if host == "" || port == 0 {
		s.mu.Lock()
		s.pendingFileRequests[friend.id] = append(s.pendingFileRequests[friend.id], pendingFileSend{
			senderID: cl.id,
			filePath: filePath,
		})
		s.mu.Unlock()
		cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Waiting for %s to get ready to receive your file...", friend.nick))
		return
	}

	s.sendFilePortToSender(cl, friend, host, port, filePath)
}

func (s *server) handleFilePortListening(receiver *client, payload json.RawMessage) {
	var meta struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
	_ = json.Unmarshal(payload, &meta)
	port, parseErr := strconv.Atoi(strings.TrimSpace(meta.Port))
	if parseErr != nil || port <= 0 || port > 65535 {
		receiver.err(fmt.Errorf("invalid file port: %q", meta.Port))
		return
	}
	host := strings.TrimSpace(meta.Host)
	if host == "" {
		receiver.err(fmt.Errorf("file port host is required"))
		return
	}

	receiver.mu.Lock()
	receiver.fileListenHost = host
	receiver.fileListenPort = port
	receiver.mu.Unlock()

	s.dispatchPendingFileRequests(receiver)
}

func (s *server) dispatchPendingFileRequests(receiver *client) {
	s.mu.Lock()
	pending := s.pendingFileRequests[receiver.id]
	delete(s.pendingFileRequests, receiver.id)
	s.mu.Unlock()
	if len(pending) == 0 {
		return
	}

	host, port := receiver.GetFriendFilePort()
	if host == "" || port == 0 {
		return
	}

	for _, req := range pending {
		s.mu.RLock()
		sender, ok := s.clients[req.senderID]
		s.mu.RUnlock()
		if !ok {
			continue
		}

		receiver.mu.RLock()
		requesterIsFriend := hasID(receiver.friends, sender.id)
		receiver.mu.RUnlock()
		if !requesterIsFriend {
			continue
		}

		s.sendFilePortToSender(sender, receiver, host, port, req.filePath)
	}
}

func (s *server) sendFilePortToSender(sender, receiver *client, host string, port int, filePath string) {
	sender.send_user_message(map[string]any{
		"host":     host,
		"port":     strconv.Itoa(port),
		"nick":     sender.nick,
		"friend":   receiver.nick,
		"filepath": filePath,
	}, FILE_PORT_LISTENING, "File port ready.")
}

func (s *server) SendRoomMessage(cl *client, roomName, message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		cl.err(fmt.Errorf("message cannot be empty"))
		return
	}
	if strings.TrimSpace(roomName) == "" {
		cl.err(fmt.Errorf("room name is required"))
		return
	}

	room_data, ok := cl.roomByName(roomName)
	if !ok {
		var err error
		room_data, err = s.resolveRoom("", roomName)
		if err != nil {
			cl.err(fmt.Errorf("you are not in room: %s", roomName))
			return
		}
		cl.mu.RLock()
		_, ok = cl.my_rooms[room_data.id]
		cl.mu.RUnlock()
		if !ok {
			cl.err(fmt.Errorf("you are not in room: %s", roomName))
			return
		}
	}

	room_data.Broadcast(cl, message)
	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Message sent to %s: %s", room_data.name, message))
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

	cl.send_user_message(map[string]any{"rooms": rooms}, DONE, formatNumberedList("these are your rooms:", names))
}

func (s *server) ListMyRoomInvites(cl *client) {
	cl.mu.RLock()
	names := make([]string, 0, len(cl.room_invites))
	for _, room_data := range cl.room_invites {
		names = append(names, room_data.name)
	}
	cl.mu.RUnlock()

	if len(names) == 0 {
		cl.send_user_message(map[string]any{}, DONE, "You have no room invites")
		return
	}

	cl.send_user_message(map[string]any{}, DONE, formatNumberedList("these are your room invites:", names))
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
