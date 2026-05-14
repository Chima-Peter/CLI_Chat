package client

import (
	"bufio"
	"net"

	"github.com/chzyer/readline"
)

func read_from_server(conn net.Conn, done chan struct{}, rl *readline.Instance) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			close(done)
			return
		}

		rl.Write([]byte(msg))
		rl.Refresh()
	}
}
