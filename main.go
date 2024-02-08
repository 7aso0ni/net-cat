package main

import (
	"log"
	"net"
	"os"

	"netcat/client"
)

var port = "8989"

func main() {
	// checking if port was passed
	if len(os.Args[1:]) != 0 {
		port = os.Args[1]
	}

	// creating a tcp connection
	listener, err := net.Listen("tcp", ":"+port) // for not it's localhost. subjected to change
	if err != nil {
		log.Fatalf("Error Listning to the tcp connection: %v", err)
	}

	defer listener.Close()

	for {
		// accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connections: %v", err)
		}

		// handle client connection in a goroutine
		go client.HandleClient(conn)
	}
}
