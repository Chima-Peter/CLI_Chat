package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func InitServer() *server {
	// create server
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) Run() {
	// listen to channel and process commands - run as goroutine
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.list_rooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) NewClient(conn net.Conn) {
	// create a new client
	log.Println("New client is connected: ", conn.RemoteAddr().String())

	new_client := client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	new_client.readInput()
}

func (s *server) nick(c *client, args []string) {
	// update user nickname
	status := c.validate_args(args, "Provide a nickname")
	if !status {
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("Nickname set to: %s", args[1]))
}

func (s *server) join(c *client, args []string) {
	// check if room exists, then add user else create new
	status := c.validate_args(args, "Provide a room name to join")
	if !status {
		return
	}

	room_name := args[1]

	r, ok := s.rooms[room_name]

	if !ok {
		r = &room{
			name:    room_name,
			members: make(map[net.Addr]*client),
		}
		s.rooms[room_name] = r
	}

	_, exists := r.members[c.conn.RemoteAddr()]
	if exists {
		c.err(fmt.Errorf("You are already a member of: %s", r.name))
		return
	}

	// add user to new room
	r.members[c.conn.RemoteAddr()] = c

	// quit current room
	s.QuitCurrentRoom(c)

	// set user room to point to current
	c.room = r

	// broadcast that new member joined
	r.Broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))

	// confirm to user that they changed rooms
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

func (s *server) list_rooms(c *client) {
	// list all rooms
	var rooms []string

	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	if len(rooms) == 0 {
		c.msg("No rooms available to join")
		return
	}

	c.msg(fmt.Sprintf("Available public rooms to join are: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	// add message to room
	status := c.validate_args(args, "Provide a valid message string")
	if !status {
		return
	}

	if c.room == nil {
		c.err(errors.New("You must join a room first!"))
		return
	}

	c.room.Broadcast(c, c.nick+": "+strings.Join(args[1:], " "))
	c.msg(fmt.Sprintf("You: %s", strings.Join(args[1:], " ")))
}

func (s *server) quit(c *client) {
	// close connection
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr())

	s.QuitCurrentRoom(c)

	c.msg("Sad to see you go")

	c.conn.Close()
}

func (s *server) QuitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.Broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
		c.msg(fmt.Sprintf("You left %s", c.room.name))
	}
}

func (s *server) CreateRoom(cl *client, room_name string) {
	if strings.TrimSpace(room_name) == "" {
		cl.err(fmt.Errorf("Provide a valid room name"))
		return
	}

	new_room := &room{
		id:         uuid.New().String(),
		name:       room_name,
		is_private: false,
		owner:      cl,
		password:   "",
		members:    make(map[net.Addr]*client),
	}

	new_room.members[cl.conn.RemoteAddr()] = cl

	s.rooms[new_room.id] = new_room

	payload := map[string]any{
		"room": room_name,
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

	room_data, ok := s.rooms[room_name]
	if !ok {
		cl.err(fmt.Errorf("No valid room found with room name: %s", room_name))
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

	room_data.DeleteRoom(cl)

	delete(s.rooms, room_data.id)
}

func (s *server) FetchRoom(cl *client, room_name string) (*room, error) {
	if strings.TrimSpace(room_name) == "" {
		return nil, fmt.Errorf("Provide room name")
	}

	room_data, ok := s.rooms[room_name]
	if !ok {
		return nil, fmt.Errorf("No valid room found with room name: %s", room_name)
	}

	return room_data, nil
}
