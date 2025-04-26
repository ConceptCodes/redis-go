package server_test

import (
	"net"
	"testing"
	"time"

	"github.com/conceptcodes/redis-go/cmd/server"
)

func TestServerStartAndAccept(t *testing.T) {
	// Arrange: Start the server in a background goroutine
	// Note: This server will continue running after the test finishes unless the program exits
	go server.Run()

	time.Sleep(100 * time.Millisecond)

	// Act: Attempt to connect to the server
	conn, err := net.DialTimeout("tcp", "localhost:6379", 2*time.Second)

	// Assert: Check if the connection was successful
	if err != nil {
		t.Fatalf("Failed to connect to server on localhost:6379: %v", err)
	}

	defer conn.Close()

	t.Log("Successfully connected to the server.")
}
