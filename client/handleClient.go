package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

// Client represents a connected client
type Client struct {
	Conn     net.Conn
	Username string
}

var (
	clients []*Client
	mu      sync.Mutex
)

var checkUsername = make(map[string]bool)

func WelcomeMsg() string {
	file, err := os.ReadFile("./WelcomeMsg.txt")
	if err != nil {
		log.Fatalf("Error Reading file: %v", err)
	}

	return string(file) + "\n"
}

func HandleClient(conn net.Conn) {
	var buffer = make([]byte, 1024)

	_, err := conn.Write([]byte(WelcomeMsg()))
	if err != nil {
		log.Fatalf("Error sending welcome message: %v", err.Error())
	}

takenUsername:
	if _, err = conn.Write([]byte("[ENTER YOUR NAME]:")); err != nil {
		log.Fatalf("Error sending welcome message: %v", err.Error())
	}

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error Reading client name: %v", err.Error())
	}

	mu.Lock()
	if _, ok := checkUsername[line]; ok {
		conn.Write([]byte("Username already taken\n"))
		goto takenUsername // this will go back to the tag and reset the operation
	} else {
		checkUsername[line] = true
	}
	mu.Unlock()

	c := &Client{Username: line, Conn: conn}
	clients = append(clients, c)

	fmt.Println(clients)

	// I have no Idea what to do with this
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
