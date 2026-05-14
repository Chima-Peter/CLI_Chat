package client

import (
	"fmt"
	"net"

	"github.com/chzyer/readline"
)

func Connect() {
	done := make(chan struct{})

	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Server currently down or not accepting new connections.")
		return
	}
	defer conn.Close()

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "> ",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	go read_from_server(conn, done, rl)

	go write_to_server(conn, rl)

	<-done
	fmt.Println("User disconnected.")
}
