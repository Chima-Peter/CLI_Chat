package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
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
			s.quit(cmd.client, cmd.args)
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
	c.nick = args[1]
	c.msg(fmt.Sprintf("Nickname set to: %s", args[1]))
}

func (s *server) join(c *client, args []string) {
	// check if room exists, then add user else create new
	room_name := args[1]

	r, ok := s.rooms[room_name]

	if !ok {
		r = &room{
			name:    room_name,
			members: make(map[net.Addr]*client),
		}
		s.rooms[room_name] = r
	}

	// add user to new room
	r.members[c.conn.RemoteAddr()] = c

	// quit current room
	s.quit_current_room(c)

	// set user room to point to current
	c.room = r

	// broadcast that new member joined
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))

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

	c.msg(fmt.Sprintf("Available rooms to join are %s: ", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	// add message to room
	if c.room == nil {
		c.err(errors.New("You must join a room first!"))
		return
	}

	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:], " "))
}

func (s *server) quit(c *client, args []string) {
	// close connection
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr())

	s.quit_current_room(c)

	c.msg("Sad to see you go")

	c.conn.Close()
}

func (s *server) quit_current_room(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
