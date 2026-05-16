package server

import (
	"fmt"
	"strings"
)

type room struct {
	id         string
	name       string
	is_private bool
	owner      *client
	password   string
	members    map[string]*client
	invites    map[string]*client
}

func (r *room) Broadcast(sender *client, msg string) {
	for id, member := range r.members {
		if id != sender.id {
			member.msg(msg)
		}
	}
}

func (r *room) SetRoomPassword(cl *client, password string) {
	trimmed_password := strings.TrimSpace(password)

	if trimmed_password == "" {
		response := cl.prepare_response(map[string]any{}, DONE, "Public room created.")

		cl.send_message(response)

		return
	}

	r.is_private = true
	r.password = trimmed_password

	response := cl.prepare_response(map[string]any{}, DONE, "Room password successfully updated.")

	cl.send_message(response)
}

func (r *room) JoinRoom(cl *client) {
	_, exists := r.members[cl.id]
	if exists {
		cl.err(fmt.Errorf("User already part of this room"))
		return
	}

	if r.is_private {
		response := cl.prepare_response(map[string]any{
			"room": r.name,
		}, GET_ROOM_PASSWORD, "Enter room password: ")

		cl.send_message(response)

		return
	}

	cl.room = r

	r.members[cl.id] = cl

	response := cl.prepare_response(map[string]any{}, DONE, "Welcome to room.")

	cl.send_message(response)
	cl.room.Broadcast(cl, fmt.Sprintf("%s joined the room", cl.nick))
}

func (r *room) JoinRoomWithPassword(cl *client, password string) {
	_, exists := r.members[cl.id]
	if exists {
		cl.err(fmt.Errorf("User already part of this room"))
		return
	}

	if r.password != password {
		cl.err(fmt.Errorf("Incorrect password."))
		return
	}

	cl.room = r

	r.members[cl.id] = cl

	response := cl.prepare_response(map[string]any{}, DONE, "Welcome to room.")

	cl.send_message(response)
	cl.room.Broadcast(cl, fmt.Sprintf("%s joined the room", cl.nick))
}

func (r *room) LeaveRoom(cl *client) {
	r.DeleteMember(cl)

	response := cl.prepare_response(map[string]any{}, DONE, "Left room.")

	cl.send_message(response)
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

	for _, client := range r.members {
		client.room = nil
		response := cl.prepare_response(map[string]any{}, DONE, "The admin has deleted this room and all members have being removed.")
		client.send_message(response)
	}

	delete(cl.my_rooms, r.id)
	delete(s.rooms, r.name)

	response := cl.prepare_response(map[string]any{}, DONE, "Deleted room.")

	cl.send_message(response)
}

func (r *room) DeleteMember(cl *client) {
	_, ok := cl.my_rooms[r.id]
	if !ok {
		cl.err(fmt.Errorf("You are not a part of this room and this error is forbidden."))
		return
	}

	delete(cl.my_rooms, r.id)
	delete(r.members, cl.id)
	if cl.room == r {
		cl.room = nil
		return
	}
}

func (r *room) FetchRoomMembers() []string {
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
	r.invites[invitee.id] = invitee
	invitee.send_message(
		invitee.prepare_response(map[string]any{}, DONE, fmt.Sprintf("Received an invite to join room: %s", r.name)),
	)

	owner.send_message(
		owner.prepare_response(map[string]any{}, DONE, fmt.Sprintf("Invite sent to %s", invitee.nick)),
	)
}

func (r *room) SeePendingRoomInvites(owner *client) {
	var invitees []string
	for _, client := range r.invites {
		invitees = append(invitees, client.nick)
	}

	owner.send_message(
		owner.prepare_response(map[string]any{}, DONE, fmt.Sprintf("These are the pending invites: %s", strings.Join(invitees, ", "))),
	)
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
		delete(r.invites, invitee.id)
		delete(invitee.room_invites, r.id)
		return
	}

	delete(r.invites, invitee.id)
	delete(invitee.room_invites, r.id)

	r.members[invitee.id] = invitee
	invitee.my_rooms[r.id] = r

	invitee.send_message(
		invitee.prepare_response(map[string]any{}, DONE, "Joined room"),
	)

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
		delete(r.invites, invitee.nick)
		delete(invitee.room_invites, r.id)
		return
	}

	delete(r.invites, invitee.nick)
	delete(invitee.room_invites, r.id)
	invitee.send_message(
		invitee.prepare_response(map[string]any{}, DONE, "You have declined this room invite"),
	)

	room_owner := r.owner
	room_owner.send_message(
		room_owner.prepare_response(map[string]any{}, DONE, fmt.Sprintf("%s has declined invitation to join room: %s", invitee.nick, r.name)),
	)
}
