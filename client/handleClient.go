package client

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Client represents a connected client
type Client struct {
	Conn     net.Conn
	Username string
	Active   bool
}

var (
	clients     []*Client
	mu          sync.Mutex
	chathistory []string
)

var checkUsername = make(map[string]bool)

func WelcomeMsg() []byte {
	file, err := os.ReadFile("./WelcomeMsg.txt")
	if err != nil {
		log.Fatalf("Error Reading file: %v", err)
	}

	return append(file, '\n')
}

func HandleClient(conn net.Conn) {
	// var buffer = make([]byte, 1024)

	_, err := conn.Write(WelcomeMsg())
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

	if name == "" || strings.Contains(name, " ") {
		conn.Write([]byte("Name shouldn't be empty or contain any spaces\n"))
		time.Sleep(2 * time.Second) // give client time to read
		goto takenUsername
	}

	mu.Lock()
	if _, ok := checkUsername[name]; ok {
		conn.Write([]byte("Username already taken\n"))
		goto takenUsername // this will go back to the tag and reset the operation
	}
	checkUsername[name] = true
	mu.Unlock()

	c := &Client{Username: name[:len(name)-1], Conn: conn, Active: true}
	if len(clients) < 10 { // limit to 10 clients per server
		clients = append(clients, c)
	} else {
		conn.Write([]byte("Too many users, get out\n")) // NEED TO TEST
		return
	}

	for _, msg := range chathistory {
		conn.Write([]byte(msg)) // Catch up the new user on previous messages
	}

	BroadcastToAllClients("Welcome " + c.Username + " to the chat!\n") // Welcome message

	// I have no Idea what to do with this
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			//log.Fatalf("Error Reading from client: %v", err.Error())
			BroadcastToAllClients(c.Username + " has disconnected...\n")
			c.Active = false
			break
		}
		BroadcastToAllClients(c.Username + ": " + line)
	}
}
