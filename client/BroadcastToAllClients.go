package client

import "fmt"

func BroadcastToAllClients(message string) {
	mu.Lock()
	defer mu.Unlock()

	for _, client := range clients {
		if _, err := client.Conn.Write([]byte(message)); err != nil {
			fmt.Println("Error writing to client")
		}
	}
}
