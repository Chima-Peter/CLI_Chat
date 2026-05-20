package server

import (
	"fmt"
	"strings"
	"sync"
)

type room struct {
	id         string
	name       string
	max_size   int
	is_private bool
	owner      *client
	password   string
	members    map[string]*client
	invites    map[string]*client
	mu         sync.RWMutex
}

func (r *room) isFull() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.max_size > 0 && len(r.members) >= r.max_size
}

func (r *room) Broadcast(sender *client, msg string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for id, member := range r.members {
		if id != sender.id {
			if sender.room.id == member.room.id {
				member.send_user_message(
					map[string]any{
						"from_user_id": sender.id,
						"from":         sender.nick,
						"message":      msg,
						"room_id":      r.id,
						"room":         r.name,
					},
					SEND_MSG,
					fmt.Sprintf("%s: %s", sender.nick, msg),
				)
			} else {
				member.send_user_message(
					map[string]any{
						"from_user_id": sender.id,
						"from":         sender.nick,
						"message":      msg,
						"room_id":      r.id,
						"room":         r.name,
					},
					SEND_MSG,
					fmt.Sprintf("[%s] %s: %s", sender.room.name, sender.nick, msg),
				)
			}
		}
	}
}

func (r *room) SetRoomPassword(cl *client, password string) {
	trimmed_password := strings.TrimSpace(password)

	if trimmed_password == "" {
		cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Public room created: %s", r.name))
		return
	}

	r.is_private = true
	r.password = trimmed_password

	cl.send_user_message(map[string]any{}, DONE, "Room password successfully updated.")
}

func (r *room) JoinRoom(cl *client) {
	r.mu.RLock()
	_, exists := r.members[cl.id]
	is_private := r.is_private
	r.mu.RUnlock()

	if cl.room != nil {
		if cl.room.id == r.id {
			cl.err(fmt.Errorf("You are currently in this room: %s", r.name))
			return
		}

		if exists && cl.room.id != r.id {
			cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Welcome to %s", cl.room.name))
			return
		}
	}

	if r.isFull() {
		cl.err(fmt.Errorf("Room is full (max %d members).", r.max_size))
		return
	}

	if is_private {
		cl.send_user_message(map[string]any{
			"room_id": r.id,
			"room":    r.name,
		}, GET_ROOM_PASSWORD, "Enter room password: ")
		return
	}

	cl.room = r
	cl.my_rooms[r.id] = r

	r.mu.Lock()
	r.members[cl.id] = cl
	r.mu.Unlock()

	cl.send_user_message(map[string]any{"room_id": r.id, "room": r.name}, DONE, fmt.Sprintf("Welcome to room: %s", r.name))
	cl.room.Broadcast(cl, fmt.Sprintf("%s joined the room", cl.nick))
}

func (r *room) JoinRoomWithPassword(cl *client, password string) {
	r.mu.RLock()
	_, exists := r.members[cl.id]
	is_password_match := r.password == password
	r.mu.RUnlock()

	if cl.room != nil {
		if cl.room.id == r.id {
			cl.err(fmt.Errorf("You are currently in this room: %s", r.name))
			return
		}

		if exists && cl.room.id != r.id {
			cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Welcome to %s", cl.room.name))
			return
		}
	}

	if !is_password_match {
		cl.err(fmt.Errorf("Incorrect password."))
		return
	}

	if r.isFull() {
		cl.err(fmt.Errorf("Room is full (max %d members).", r.max_size))
		return
	}

	cl.room = r
	cl.my_rooms[r.id] = r

	r.mu.Lock()
	r.members[cl.id] = cl
	r.mu.Unlock()

	cl.send_user_message(map[string]any{"room_id": r.id, "room": r.name}, DONE, fmt.Sprintf("Welcome to %s", r.name))
	cl.room.Broadcast(cl, fmt.Sprintf("%s joined the room", cl.nick))
}

func (r *room) EditRoom(cl *client, s *server, newName string, maxSize *int) {
	if r.owner.id != cl.id {
		cl.err(fmt.Errorf("Only the room owner can edit this room."))
		return
	}

	newName = strings.TrimSpace(newName)

	r.mu.Lock()
	if maxSize != nil {
		if *maxSize < 0 {
			r.mu.Unlock()
			cl.err(fmt.Errorf("Max size cannot be negative."))
			return
		}
		if *maxSize > 0 && len(r.members) > *maxSize {
			r.mu.Unlock()
			cl.err(fmt.Errorf("Cannot set max size below current member count (%d).", len(r.members)))
			return
		}
		r.max_size = *maxSize
	}
	oldName := r.name
	r.mu.Unlock()

	if newName != "" && newName != oldName {
		if s.isRoomNameTaken(newName, r.id) {
			cl.err(fmt.Errorf("Room name %s is already in use.", newName))
			return
		}

		r.mu.Lock()
		r.name = newName
		r.mu.Unlock()
	}

	r.mu.RLock()
	payload := map[string]any{
		"room_id":  r.id,
		"room":     r.name,
		"max_size": r.max_size,
	}
	r.mu.RUnlock()

	cl.send_user_message(payload, DONE, "Room settings updated.")
}

func (r *room) LeaveRoom(cl *client) {
	r.DeleteMember(cl)

	cl.send_user_message(map[string]any{}, DONE, "Left room.")
	r.Broadcast(cl, fmt.Sprintf("%s left the room", cl.nick))
}

func (r *room) DeleteRoom(cl *client, s *server) {
	_, ok := cl.my_rooms[r.id]
	if !ok {
		cl.err(fmt.Errorf("You are not a part of this room and this error is forbidden."))
	}

	if r.owner.id != cl.id {
		cl.err(fmt.Errorf("You cannot perform this action as you're not this room's admin."))
	}

	r.mu.RLock()
	members_copy := make(map[string]*client)
	for id, member := range r.members {
		members_copy[id] = member
	}
	r.mu.RUnlock()

	for _, client := range members_copy {
		client.room = nil
		client.send_user_message(map[string]any{}, DONE, "The admin has deleted this room and all members have being removed.")
	}

	delete(cl.my_rooms, r.id)
	s.mu.Lock()
	delete(s.rooms, r.id)
	s.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, "Deleted room.")
}

func (r *room) DeleteMember(cl *client) {
	_, ok := cl.my_rooms[r.id]
	if !ok {
		cl.err(fmt.Errorf("You are not a part of this room and this error is forbidden."))
		return
	}

	delete(cl.my_rooms, r.id)
	r.mu.Lock()
	delete(r.members, cl.id)
	r.mu.Unlock()
	if cl.room == r {
		cl.room = nil
		return
	}
}

func (r *room) FetchRoomMembers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var members []string

	if len(r.members) == 0 {
		return members
	}

	for _, client := range r.members {
		members = append(members, client.nick)
	}

	return members
}

func (r *room) SendRoomInvite(invitee *client, owner *client) {
	invitee.room_invites[r.id] = r
	r.mu.Lock()
	r.invites[invitee.id] = invitee
	r.mu.Unlock()
	invitee.send_user_message(map[string]any{
		"room_id": r.id,
		"room":    r.name,
	}, DONE, fmt.Sprintf("Received an invite to join room: %s", r.name))
	owner.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Invite sent to %s", invitee.nick))
}

func (r *room) SeePendingRoomInvites(owner *client) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var invitees []string
	for _, client := range r.invites {
		invitees = append(invitees, client.nick)
	}

	owner.send_user_message(map[string]any{}, DONE, fmt.Sprintf("These are the pending invites: %s", strings.Join(invitees, ", ")))
}

func (r *room) AcceptRoomInvite(invitee *client) {
	_, exists := invitee.room_invites[r.id]
	if !exists {
		invitee.err(fmt.Errorf("You have not been invited to this room"))
		return
	}

	_, ok := invitee.my_rooms[r.id]
	if ok {
		invitee.err(fmt.Errorf("You are already a part of this room"))
		r.mu.Lock()
		delete(r.invites, invitee.id)
		r.mu.Unlock()
		delete(invitee.room_invites, r.id)
		return
	}

	if r.isFull() {
		invitee.err(fmt.Errorf("Room is full (max %d members).", r.max_size))
		return
	}

	r.mu.Lock()
	delete(r.invites, invitee.id)
	r.members[invitee.id] = invitee
	r.mu.Unlock()
	delete(invitee.room_invites, r.id)

	invitee.my_rooms[r.id] = r
	invitee.room = r

	invitee.send_user_message(map[string]any{"room_id": r.id, "room": r.name}, DONE, "Joined room")

	r.Broadcast(invitee, fmt.Sprintf("%s just joined the room.", invitee.nick))
}

func (r *room) DeclineRoomInvite(invitee *client) {
	_, exists := invitee.room_invites[r.id]
	if !exists {
		invitee.err(fmt.Errorf("You have not been invited to this room"))
		return
	}

	_, ok := invitee.my_rooms[r.id]
	if ok {
		invitee.err(fmt.Errorf("You are already a part of this room"))
		r.mu.Lock()
		delete(r.invites, invitee.id)
		r.mu.Unlock()
		delete(invitee.room_invites, r.id)
		return
	}

	r.mu.Lock()
	delete(r.invites, invitee.id)
	r.mu.Unlock()
	delete(invitee.room_invites, r.id)
	invitee.send_user_message(map[string]any{}, DONE, "You have declined this room invite")

	room_owner := r.owner
	room_owner.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s has declined invitation to join room: %s", invitee.nick, r.name))
}
