package server_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/conceptcodes/redis-go/cmd/server"
	"github.com/conceptcodes/redis-go/internal/constants"
)

func TestServerStartAndAccept(t *testing.T) {
	// Arrange: Start the server in a background goroutine
	go server.Start()

	t.Cleanup(server.Shutdown)

	time.Sleep(100 * time.Millisecond)

	// Act: Attempt to connect to the server
	url := fmt.Sprintf("localhost:%d", constants.Port)
	conn, err := net.DialTimeout("tcp", url, 2*time.Second)

	// Assert: Check if the connection was successful
	if err != nil {
		t.Fatalf("Failed to connect to server on localhost:5555: %v", err)
	}

	defer conn.Close()

	t.Log("Successfully started the server.")
}

func TestServerShutdown(t *testing.T) {
	// Arrange: Start the server in a background goroutine
	go server.Start()

	t.Cleanup(server.Shutdown)

	url := fmt.Sprintf("localhost:%d", constants.Port)
	time.Sleep(100 * time.Millisecond)

	// Act: Shutdown the server
	server.Shutdown()

	// Assert: Check if the server is no longer accepting connections
	if _, err := net.DialTimeout("tcp", url, 2*time.Second); err == nil {
		t.Fatal("Server should not accept connections after shutdown")
	}

	t.Log("Server successfully shut down.")
}
