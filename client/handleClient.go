package client

import (
	"log"
	"net"
	"os"
)

// Client represents a connected client
type Client struct {
	Conn     net.Conn
	Username string
}

func WelcomeMsg() string {
	file, err := os.ReadFile("./WelcomeMsg.txt")
	if err != nil {
		log.Fatalf("Error Reading file: %v", err)
	}

	return string(file) + "\n"
}

func HandleClient(conn net.Conn) {
	penguin := WelcomeMsg()
	_, err := conn.Write([]byte(penguin + "[ENTER YOUR NAME]:"))
	if err != nil {
		log.Fatalf("Error sending welcome message: %v", err.Error())
	}

	buffer := make([]byte, 1024)

	// I have no Idea what is this
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			log.Fatalf("something went wrong with reading from client: %v", err)
		}

		// printing what is written to the client
		_, err = conn.Write(buffer)
		if err != nil {
			log.Fatal("Error Writing data to server")
		}
	}

}
