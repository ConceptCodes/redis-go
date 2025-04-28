package server

import (
	"fmt"
	"log"
	"net"

	"github.com/conceptcodes/redis-go/internal/constants"
)

func Start() {
	port := fmt.Sprintf(":%d", constants.Port)
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to bind to port %s: %v", port, err)
	}
	defer l.Close()
	log.Printf("Listening on port %s", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	log.Printf("Serving %s", c.RemoteAddr().String())
	// TODO: read from the connection, parse commands,
	log.Printf("Connection closed for %s", c.RemoteAddr().String())
}
