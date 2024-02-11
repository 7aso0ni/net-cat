package client

import (
	"log"
	"net"
)

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("Error getting local ip: %v", err)
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
