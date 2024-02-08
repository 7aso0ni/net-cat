package client

import (
	"log"
	"net"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

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
