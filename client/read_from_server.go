package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/chzyer/readline"
)

func read_from_server(conn net.Conn, rl *readline.Instance) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return
		}

		rl.Clean()
		fmt.Fprintf(rl.Stdout(), "> %s", msg)
		rl.Refresh()
	}
}
