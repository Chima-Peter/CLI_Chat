package client

import (
	"fmt"
	"net"

	"github.com/chzyer/readline"
)

func printBootstrap() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║         CLI Chat Client                ║")
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Println("║  Commands:                             ║")
	fmt.Println("║    /nick <name>  - Set your nickname   ║")
	fmt.Println("║    /join <room>  - Join a chat room    ║")
	fmt.Println("║    /rooms        - List all rooms      ║")
	fmt.Println("║    /msg <text>   - Send a message      ║")
	fmt.Println("║    /quit         - Disconnect          ║")
	fmt.Println("╚════════════════════════════════════════╝")
}

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

	go write_to_server(conn, rl)
	read_from_server(conn, rl)

	fmt.Println("\n> Disconnected.")
}
