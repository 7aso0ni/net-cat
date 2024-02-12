package client

import (
	"fmt"
	"os"
	"time"
)

func Broadcast(message string, exclude *Client, saveMessage bool) {
	mu.Lock()
	defer mu.Unlock()
	if saveMessage { // Write to log file and chatHistory array
		logFile, err := os.OpenFile("logs.txt",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
		if _, err := logFile.WriteString(message); err != nil {
			fmt.Println("Error writing to log file:" + message)
		}
		logFile.Close()
		chathistory = append(chathistory, message)
	}

	for _, client := range clients {
		if client != exclude {
			// Send recieved message to each client
			if _, err := client.Conn.Write([]byte("\n" + message)); err != nil {
				fmt.Println("Error writing to client:", err)
			}
		}
		// Write client specific header
		if _, err := client.Conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Username + "]:")); err != nil {
			fmt.Println("Error writing to client:", err)
		}
	}
}

//func BroadcastToAllClients(message string) {
//	mu.Lock()
//	defer mu.Unlock()
//	// Add message to full file log
//	logFile, err := os.OpenFile("logs.txt",
//		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	defer logFile.Close()
//	logFile.WriteString(message)
//	// Add message to local log
//	chathistory = append(chathistory, message)
//	for _, client := range clients {
//		if _, err := client.Conn.Write([]byte(message)); err != nil {
//			fmt.Println("Error writing to client: " + client.Username)
//		}
//	}
//}
//
//func BroadcastMessageToOthers(message string, exclude *Client) {
//	mu.Lock()
//	defer mu.Unlock()
//	for _, client := range clients {
//		if client != exclude {
//			_, err := client.Conn.Write([]byte(message))
//			if err != nil {
//				fmt.Println("Error writing to client:", err)
//			}
//		}
//	}
//}
