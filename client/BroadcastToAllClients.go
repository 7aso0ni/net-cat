package client

import "fmt"

func BroadcastToAllClients(message string) {
	chathistory = append(chathistory, message) // Add to universal log
	mu.Lock()
	defer mu.Unlock()

	for _, client := range clients {
		if client.Active {
			if _, err := client.Conn.Write([]byte(message)); err != nil {
				fmt.Println("Error writing to client: " + client.Username)
			}
		}
	}
}
