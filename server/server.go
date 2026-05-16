package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

type server struct {
	rooms    map[string]*room
	commands chan command
	clients  map[string]*client
}

func InitServer() *server {
	// create server
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
		clients:  make(map[string]*client),
	}
}

func (s *server) Run() {
	// listen to channel and process commands - run as goroutine
	for cmd := range s.commands {
		switch cmd.id {
		}
	}
}

func (s *server) NewClient(conn net.Conn) {
	// create a new client
	log.Println("New client is connected: ", conn.RemoteAddr().String())

	new_client := client{
		id:           uuid.New().String(),
		conn:         conn,
		nick:         "anonymous",
		my_rooms:     make(map[string]*room),
		room_invites: make(map[string]*room),
		commands:     s.commands,
	}

	s.clients[new_client.id] = &new_client

	new_client.readInput()
}

func (s *server) CreateRoom(cl *client, room_name string) {
	trimmed_name := strings.TrimSpace(room_name)
	if trimmed_name == "" {
		cl.err(fmt.Errorf("Provide a valid room name"))
		return
	}

	_, exists := s.rooms[trimmed_name]
	if exists {
		cl.err(fmt.Errorf("%s has already been used. Select another name.", room_name))
	}

	new_room := &room{
		id:         uuid.New().String(),
		name:       trimmed_name,
		is_private: false,
		owner:      cl,
		password:   "",
		members:    make(map[string]*client),
		invites:    make(map[string]*client),
	}

	new_room.members[cl.id] = cl

	s.rooms[new_room.name] = new_room

	payload := map[string]any{
		"room": trimmed_name,
	}

	response := cl.prepare_response(payload, SET_ROOM_PASSWORD, "Set room password. Leave black for public room: ")

	cl.send_message(response)
}

func (s *server) SetRoomPassword(cl *client, room_name string, password string) {
	room_data, err := s.FetchRoom(cl, room_name)

	if err != nil {
		cl.err(err)
	}

	if room_data.owner.id != cl.id {
		cl.err(fmt.Errorf("You don't have the right to update this room password."))
		return
	}

	room_data.SetRoomPassword(cl, password)
}

func (s *server) JoinRoom(cl *client, room_name string) {
	room_data, err := s.FetchRoom(cl, room_name)

	if err != nil {
		cl.err(err)
	}

	room_data.JoinRoom(cl)
}

func (s *server) JoinRoomWithPassword(cl *client, room_name string, password string) {
	trimmed_password := strings.TrimSpace(password)

	if strings.TrimSpace(room_name) == "" {
		cl.err(fmt.Errorf("Provide room name"))
		return
	}

	if trimmed_password == "" {
		cl.err(fmt.Errorf("Enter password"))

		return
	}

	room_data, err := s.GetRoomByName(room_name)
	if err != nil {
		cl.err(err)
		return
	}

	room_data.JoinRoomWithPassword(cl, password)
}

func (s *server) LeaveRoom(cl *client, room_name string) {
	room_data, err := s.FetchRoom(cl, room_name)

	if err != nil {
		cl.err(err)
	}

	room_data.LeaveRoom(cl)
}

func (s *server) DeleteRoom(cl *client, room_name string) {
	room_data, err := s.FetchRoom(cl, room_name)

	if err != nil {
		cl.err(err)
	}

	room_data.DeleteRoom(cl, s)
}

func (s *server) DeleteRoomMember(cl *client, room_name string, client_name string) {
	room_data, err := s.FetchRoom(cl, room_name)
	if err != nil {
		cl.err(err)
	}

	if room_data.owner.id != cl.id {
		cl.err(fmt.Errorf("This action is reserved for only the admin!"))
	}

	member, err := s.GetClientByNick(client_name)
	if err != nil {
		cl.err(err)
	}

	_, ok := room_data.members[member.id]
	if !ok {
		cl.err(fmt.Errorf("%s is not a member of room: %s", client_name, room_name))
	}

	room_data.DeleteMember(member)

	delete(s.rooms, room_data.name)

	response := cl.prepare_response(map[string]any{}, DONE, fmt.Sprintf("%s has been deleted from this room.", client_name))

	member.send_message(
		member.prepare_response(map[string]any{}, DONE, fmt.Sprintf("Admin deleted you from room: %s", room_name)),
	)

	room_data.Broadcast(cl, fmt.Sprintf("%s has been deleted from this room", client_name))

	cl.send_message(response)
}

func (s *server) GetRoomMembers(cl *client, room_name string) {
	room_data, err := s.FetchRoom(cl, room_name)

	if err != nil {
		cl.err(err)
	}

	_, exists := room_data.members[cl.id]
	if !exists {
		cl.err(fmt.Errorf("Only members of a room can see it's members"))
	}

	members := room_data.FetchRoomMembers()

	cl.send_message(
		cl.prepare_response(map[string]any{}, DONE, fmt.Sprintf("Members of this room are: %s", strings.Join(members, ", "))),
	)
}

func (s *server) FetchRoom(cl *client, room_name string) (*room, error) {
	return s.GetRoomByName(room_name)
}

// GetClientByNick returns a client by their nickname
func (s *server) GetClientByNick(nick string) (*client, error) {
	for _, client := range s.clients {
		if client.nick == nick {
			return client, nil
		}
	}
	return nil, fmt.Errorf("Client not found with nick: %s", nick)
}

// GetClientIDByNick returns a client's ID by their nickname
func (s *server) GetClientIDByNick(nick string) (string, error) {
	client, err := s.GetClientByNick(nick)
	if err != nil {
		return "", err
	}
	return client.id, nil
}

// GetRoomByName returns a room by its name
func (s *server) GetRoomByName(room_name string) (*room, error) {
	if strings.TrimSpace(room_name) == "" {
		return nil, fmt.Errorf("Provide room name")
	}

	room_data, ok := s.rooms[room_name]
	if !ok {
		return nil, fmt.Errorf("No valid room found with room name: %s", room_name)
	}

	return room_data, nil
}

// GetRoomIDByName returns a room's ID by its name
func (s *server) GetRoomIDByName(room_name string) (string, error) {
	room_data, err := s.GetRoomByName(room_name)
	if err != nil {
		return "", err
	}
	return room_data.id, nil
}

func (s *server) SendRoomInvite(cl *client, room_name string, client_name string) {
	room_data, err := s.FetchRoom(cl, room_name)
	if err != nil {
		cl.err(err)
		return
	}

	if room_data.owner.id != cl.id {
		cl.err(fmt.Errorf("This action is reserved for only the admin!"))
		return
	}

	new_member, err := s.GetClientByNick(client_name)
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

func (s *server) SeePendingRoomInvites(cl *client, room_name string) {
	room_data, err := s.FetchRoom(cl, room_name)
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

func (s *server) AcceptRoomInvite(invitee *client, room_name string) {
	room_data, err := s.FetchRoom(invitee, room_name)
	if err != nil {
		invitee.err(err)
		return
	}

	room_data.AcceptRoomInvite(invitee)
}

func (s *server) DeclineRoomInvite(invitee *client, room_name string) {
	room_data, err := s.FetchRoom(invitee, room_name)
	if err != nil {
		invitee.err(err)
		return
	}

	room_data.AcceptRoomInvite(invitee)
}

func (s *server) ListPublicRooms(cl *client) {
	var response []string

	for _, room := range s.rooms {
		if !room.is_private {
			response = append(response, room.name)
		}
	}

	cl.send_message(
		cl.prepare_response(map[string]any{}, DONE, fmt.Sprintf("Publicly available rooms are: %s", strings.Join(response, ", "))),
	)
}

func (s *server) ListMyRooms(cl *client) {
	var response []string

	for _, room := range cl.my_rooms {
		response = append(response, room.name)
	}

	cl.send_message(
		cl.prepare_response(map[string]any{}, DONE, fmt.Sprintf("You are a member of: %s", strings.Join(response, ", "))),
	)
}
