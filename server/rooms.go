package server

import (
	"fmt"
	"net"
	"slices"
	"strings"
)

type room struct {
	id         string
	name       string
	is_private bool
	owner      *client
	password   string
	members    map[net.Addr]*client
}

func (r *room) Broadcast(sender *client, msg string) {
	for addr, member := range r.members {
		if addr != sender.conn.RemoteAddr() {
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
	if r.is_private {
		response := cl.prepare_response(map[string]any{
			"room": r.name,
		}, GET_ROOM_PASSWORD, "Enter room password: ")

		cl.send_message(response)

		return
	}

	cl.room = r

	response := cl.prepare_response(map[string]any{}, DONE, "Welcome to room.")

	cl.send_message(response)
	cl.room.Broadcast(cl, fmt.Sprintf("%s joined the room", cl.nick))
}

func (r *room) JoinRoomWithPassword(cl *client, password string) {
	if r.password != password {
		cl.err(fmt.Errorf("Incorrect password."))
		return
	}

	cl.room = r

	response := cl.prepare_response(map[string]any{}, DONE, "Welcome to room.")

	cl.send_message(response)
	cl.room.Broadcast(cl, fmt.Sprintf("%s joined the room", cl.nick))
}

func (r *room) LeaveRoom(cl *client) {
	idx := slices.IndexFunc(cl.my_rooms, func(room room) bool {
		return room.name == r.name
	})

	if idx == -1 {
		cl.err(fmt.Errorf("You are not a part of this room"))
		return
	}

	r.DeleteMember(r, cl, idx)

	response := cl.prepare_response(map[string]any{}, DONE, "Left room.")

	cl.send_message(response)
	r.Broadcast(cl, fmt.Sprintf("%s left the room", cl.nick))
}

func (r *room) DeleteRoom(cl *client) {
	idx := slices.IndexFunc(cl.my_rooms, func(room room) bool {
		return room.name == r.name
	})

	if idx == -1 {
		cl.err(fmt.Errorf("You are not a part of this room"))
	}

	if r.owner.id != cl.id {
		cl.err(fmt.Errorf("You cannot perform this action as you're not this room's admin."))
	}

	for _, client := range r.members {
		client.room = nil
		response := cl.prepare_response(map[string]any{}, DONE, "The admin has deleted this room and all members have being removed.")
		client.send_message(response)
	}

	response := cl.prepare_response(map[string]any{}, DONE, "Deleted room.")

	cl.send_message(response)
}

func (r *room) DeleteMember(room_data *room, cl *client, idx int) {
	cl.my_rooms = append(cl.my_rooms[:idx], cl.my_rooms[idx+1:]...)
	delete(room_data.members, cl.conn.RemoteAddr())
	if cl.room == room_data {
		cl.room = nil
	}
}
