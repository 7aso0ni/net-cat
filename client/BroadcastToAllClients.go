package client

import (
	"fmt"
	"netcat/ui"
	"os"
	"time"
)

func Broadcast(message string, exclude *Client, saveMessage bool) {
	mu.Lock()
	defer mu.Unlock()
	if saveMessage { // Write to log file and chatHistory array
		ui.AddMessage(message)
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
			if _, err := client.Conn.Write([]byte("\n" + message + headerStr(client.Username))); err != nil {
				fmt.Println("Error writing to client:", err)
			}
		} else {
			// Write client specific header
			if _, err := client.Conn.Write([]byte(headerStr(client.Username))); err != nil {
				fmt.Println("Error writing to client:", err)
			}
		}
	}
}

func headerStr(username string) string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + username + "]:"
}
