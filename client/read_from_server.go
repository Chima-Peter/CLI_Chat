package client

import (
	"bufio"
	"fmt"
	"net"
)

func read_from_server(conn net.Conn, done chan struct{}) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			close(done)
			return
		}

		fmt.Printf("%s", msg)
	}
}
