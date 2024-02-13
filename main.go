package main

import (
	"fmt"
	"net"
	"netcat/client"
	"netcat/ui"
	"os"
)

var port = "8989"
var ip net.IP

func main() {
	// checking if port was passed
	if len(os.Args[1:]) != 0 {
		port = os.Args[1]
	}

	// creating a tcp connection
	ip = client.GetLocalIP()
	listener, err := net.Listen("tcp", ip.String()+":"+port) // listen on local ip
	if err != nil {
		conn, err := net.Dial("tcp", ":"+port)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		defer conn.Close()
		//log.Fatalf("Error Listning to the tcp connection: %v", err)

	}

	defer listener.Close()

	go ui.OpenUI()
	ui.Header = fmt.Sprintf("Listening on host %v:%v", ip, port)

	for {
		// accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connections: %s\n", err.Error())
			continue // if error encounterd try again
		}
		// handle client connection in a goroutine
		go client.HandleClient(conn)
	}
}
