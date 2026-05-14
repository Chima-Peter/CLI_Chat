package client

import (
	"fmt"
	"net"
)

func Connect() {
	done := make(chan struct{})

	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Server currently down or not accepting new connections.")
		return
	}
	defer conn.Close()

	go read_from_server(conn, done)

	go write_to_server(conn)

	<-done
	fmt.Println("> User disconnected.")
}
