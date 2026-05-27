package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/chima/CLI_Chat/server"
)

func printBootstrap() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║         CLI Chat Server                ║")
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Println("║  Server Info:                          ║")
	fmt.Println("║    Port: 8888                          ║")
	fmt.Println("║    Protocol: TCP                       ║")
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Println("║  Client Commands:                      ║")
	fmt.Println("║    /nick <name>  - Set nickname        ║")
	fmt.Println("║    /join <room>  - Join a room         ║")
	fmt.Println("║    /rooms        - List rooms          ║")
	fmt.Println("║    /msg <text>   - Send message        ║")
	fmt.Println("║    /quit         - Disconnect          ║")
	fmt.Println("╚════════════════════════════════════════╝")
}

func main() {
	printBootstrap()

	srv := server.InitServer()

	// start listening on port 8888
	listener, err := tls.Listen("tcp", ":8888", &tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Fatal("Unable to start server:", err.Error())
	}

	defer listener.Close()

	log.Println("Server is running on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Unable to accept connection:", err.Error())
			continue
		}

		go srv.HandleConn(conn)
	}
}
