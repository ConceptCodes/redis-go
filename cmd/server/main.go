package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/conceptcodes/redis-go/internal/constants"
	"github.com/conceptcodes/redis-go/internal/resp"
)

var l net.Listener

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

func Shutdown() {
	log.Println("Server shutting down...")
	if l != nil {
		l.Close()
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	log.Printf("Serving %s", c.RemoteAddr().String())

	p := resp.NewParser(c)
	for {
		cmd, err := p.Parse()
		if err != nil {
			if err == io.EOF {
				log.Printf("Client %s disconnected", c.RemoteAddr().String())
				break
			}
			log.Printf("Error parsing command from %s: %v", c.RemoteAddr().String(), err)
			c.Write([]byte(err.Error()))
		}
		if cmd == nil {
			log.Printf("Received nil command from %s", c.RemoteAddr().String())
			continue
		}
		log.Printf("Received command from %s: %s", c.RemoteAddr().String(), cmd)
		// TODO: Process the command and send a response
	}
}
