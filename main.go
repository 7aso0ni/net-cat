package main

import (
	"fmt"
	"log"
	"net"
	"netcat/client"
	"os"
)

var port = "8989"

func main() {
	// checking if port was passed
	if len(os.Args[1:]) != 0 {
		port = os.Args[1]
	}

	// creating a tcp connection
	listener, err := net.Listen("tcp", ":"+port) // for now it's localhost. subjected to change
	if err != nil {
		conn, err := net.Dial("tcp", ":"+port)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer conn.Close()
		//log.Fatalf("Error Listning to the tcp connection: %v", err)

	}

	defer listener.Close()
	log.Printf("Listening on port %v", port)

	for {
		// accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connections: %s\n", err.Error())
			continue // if error encounterd try again
		}

		// handle client connection in a goroutine
		go client.HandleClient(conn)
	}
}
