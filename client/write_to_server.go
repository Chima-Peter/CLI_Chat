package client

import (
	"bufio"
	"net"
	"os"
)

func write_to_server(conn net.Conn) {
	input := bufio.NewReader(os.Stdin)

	for {
		msg, err := input.ReadString('\n')
		if err != nil {
			return
		}

		conn.Write([]byte(msg))
	}
}
