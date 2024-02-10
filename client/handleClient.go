package client

import (
	"bufio"
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
	// var buffer = make([]byte, 1024)

	_, err := conn.Write([]byte(WelcomeMsg()))
	if err != nil {
		log.Fatalf("Error sending welcome message: %v", err.Error())
	}

takenUsername:
	if _, err = conn.Write([]byte("[ENTER YOUR NAME]:")); err != nil {
		log.Fatalf("Error sending welcome message: %v", err.Error())
	}

	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n') //read the entire line from the client
	if err != nil {
		log.Fatalf("Error Reading client name: %v", err.Error())
	}

	mu.Lock()
	if _, ok := checkUsername[name]; ok {
		conn.Write([]byte("Username already taken\n"))
		goto takenUsername // this will go back to the tag and reset the operation
	} else {
		checkUsername[name] = true
	}
	mu.Unlock()

	c := &Client{Username: name, Conn: conn}
	if len(clients) < 10 { // limit to 10 clients per server
		clients = append(clients, c)
	}

	for _, client := range clients {
		if client.Username != name {
			client.Conn.Write([]byte(name + " has joined the chat..."))
		}
	}

	// I have no Idea what to do with this
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error Reading from client: %v", err.Error())
		}
		BroadcastToAllClients(line)
	}

}
