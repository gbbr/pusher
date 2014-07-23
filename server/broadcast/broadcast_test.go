package broadcast

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http/httptest"
	"testing"
)

const (
	ORIGIN       = "http://localhost"
	SOCKET_PATH  = "/test"
	FULL_ADDRESS = "ws://%s" + SOCKET_PATH
	TEST_MESSAGE = "Ping Pong"
)

func TestBroadcast(t *testing.T) {
	var (
		addr    string
		clients [5]*websocket.Conn
		err     error
		msg     string
	)

	// Start broadcast on set path
	server := New()
	server.Start(SOCKET_PATH)

	// Test server
	ts := httptest.NewServer(nil)
	addr = ts.Listener.Addr().String()
	fmt.Println("Test server listening on ", addr)

	// Set up 5 clients
	for i := 0; i < 5; i++ {
		clients[i], err = websocket.Dial(fmt.Sprintf(FULL_ADDRESS, addr), "", ORIGIN)
		if err != nil {
			fmt.Println("Failed connecting")
		}
	}

	// Send message from client 1 and test broadcast
	websocket.Message.Send(clients[0], TEST_MESSAGE)
	for i := 0; i < 5; i++ {
		websocket.Message.Receive(clients[i], &msg)
		if msg != TEST_MESSAGE {
			fmt.Println("Message not received or invalid")
			t.Fail()
		}
	}
}
