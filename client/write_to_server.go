package client

import (
	"fmt"
	"net"

	"github.com/chzyer/readline"
)

func write_to_server(conn net.Conn, rl *readline.Instance) {
	for {
		msg, err := rl.Readline()
		if err != nil {
			rl.Clean()
			conn.Write([]byte("/quit\n"))
			conn.Close()
			return
		}

		fmt.Fprint(rl.Stdout(), "\033[A\033[2K")
		conn.Write([]byte(msg + "\n"))
	}
}
