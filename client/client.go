package client

import (
	"fmt"
	"net"

	"github.com/chzyer/readline"
)

func Connect() {
	printBootstrap()

	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Server currently down or not accepting new connections.")
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server!")

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "\033[5m▌\033[0m ",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	runSession(conn, rl)

	fmt.Println("\n> Disconnected.")
}
