package main

import (
	"log"
	"net"

	"github.com/chima/CLI_Chat/server"
)

func main() {
	// initialize server
	server := server.InitServer()

	// listen to channel in background and process commands
	go server.Run()

	// start listening on port 8888
	listener, err := net.Listen("tcp", ":8888")

	if err != nil {
		log.Fatal("Unable to start server:", err.Error())
	}

	defer listener.Close()

	log.Println("Started server on :8888")

	for {
		// accept connections infinitely on port 8888
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Unable to accept connection:", err.Error())
			continue
		}

		go server.NewClient(conn)
	}
}
