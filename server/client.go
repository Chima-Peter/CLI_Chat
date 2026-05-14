package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	// read input from TCP connection using bufio.NewReader that accepts a readable stream as far as it implements io.Reader, then we add a delimiter, process the input. Doing so in an infinite loop
	reader := bufio.NewReader(c.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")

		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("Unknown command: %s", cmd))
		}
	}
}

func (client *client) err(err error) {
	client.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (client *client) msg(msg string) {
	client.conn.Write([]byte("> " + msg + "\n"))
}
