package client

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/chzyer/readline"
)

func Connect() {
	printBootstrap()

	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		serverUrl = "localhost:" + port
	}

	conn, err := tls.Dial("tcp", serverUrl, &tls.Config{
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
