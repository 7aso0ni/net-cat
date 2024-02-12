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
	UID      int
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

	name = name[:len(name)-1]
	if name == "" || strings.Contains(name, " ") {
		conn.Write([]byte("Name shouldn't be empty or contain any spaces\n"))
		time.Sleep(1 * time.Second) // give client time to read
		goto takenUsername          // this will go back to the tag and reset the operation
	}

	// making the username unique
	if checkUsername[strings.ToLower(name)] {
		conn.Write([]byte("Username already taken\n"))
		time.Sleep(1 * time.Second)
		goto takenUsername
	}
	checkUsername[strings.ToLower(name)] = true

	c := &Client{Username: name, Conn: conn, UID: len(clients)}
	if len(clients) < 10 { // limit to 10 clients per server
		clients = append(clients, c)
	} else {
		conn.Write([]byte("Too many users, get out\n"))
		conn.Close()
		return
	}

	// shared resourses should be locked to prevent errors
	mu.Lock()
	for _, msg := range chathistory {
		conn.Write([]byte(msg)) // Catch up the new user on previous messages
	}
	mu.Unlock()

	Broadcast(c.Username+" has joined out chat!\n", c, false) // Welcome message

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			Broadcast(c.Username+" has left our chat...\n", c, false) // Exit message
			DeleteClient(c.UID)
			break
		}
		if line != "\n" {
			Broadcast("["+time.Now().Format("2006-01-02 15:04:05")+"]["+c.Username+"]:"+line, c, true)
		} else {
			conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + c.Username + "]:"))
		}
	}
}

func DeleteClient(UID int) {
	mu.Lock()
	defer mu.Unlock()
	delete(checkUsername, clients[UID].Username)        // Forget Username
	clients = append(clients[:UID], clients[UID+1:]...) // Remove Username
	for UID < len(clients) {                            //Update UIDs
		clients[UID].UID = UID
		UID++
	}
}
