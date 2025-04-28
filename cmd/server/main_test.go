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

	time.Sleep(100 * time.Millisecond)

	// Act: Attempt to connect to the server
	url := fmt.Sprintf("localhost:%d", constants.Port)
	conn, err := net.DialTimeout("tcp", url, 2*time.Second)

	// Assert: Check if the connection was successful
	if err != nil {
		t.Fatalf("Failed to connect to server on localhost:5555: %v", err)
	}

	defer conn.Close()

	t.Log("Successfully connected to the server.")
}
