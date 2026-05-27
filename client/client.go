package client

import (
	"crypto/tls"
	"fmt"

	"github.com/chzyer/readline"
)

func Connect() {
	printBootstrap()

	conn, err := tls.Dial("tcp", "localhost:8888", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println("> Server currently down or not accepting new connections.")
		return
	}
	defer conn.Close()

	fmt.Println("> Connected to server!")

	rl, err := readline.NewEx(&readline.Config{
		Prompt: defaultPrompt,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	runSession(conn, rl)

	fmt.Printf("\n> Disconnected.")
}
