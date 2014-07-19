package broadcast

import (
	"code.google.com/p/go.net/websocket"
	"io"
	"net/http"
)

// Initializes a new server
func New() *Server {
	return &Server{make([]Client, 0, 10), make(chan string)}
}

/*
CLIENT
Listens to incoming messages and sends them to socket
*/
type Client struct {
	connection *websocket.Conn
	terminate  chan bool
}

// Send back via socket when incoming message is available
func (client *Client) Listen(incoming chan string) {
LISTENER_LOOP:
	for {
		select {
		case msg := <-incoming:
			websocket.Message.Send(client.connection, msg)
		case <-client.terminate:
			break LISTENER_LOOP
		}
	}
}

// Terminate connection and trigger breaking of listener
func (client *Client) Terminate() {
	client.connection.Close()
	client.terminate <- true
}

/*
SERVER
Receives connections, creates new clients and feeds incoming messages
to communication channel
*/
type Server struct {
	connections []Client
	broadcast   chan string
}

// Listen to WebSocket connections and register clients to
// communication channels
func (server *Server) Start(path string) {
	onConnected := func(ws *websocket.Conn) {
		defer ws.Close()
		server.AddClient(Client{ws, make(chan bool)})
	}

	http.Handle(path, websocket.Handler(onConnected))
}

// Receive message from connection and send to communication channel
func (server *Server) AddClient(client Client) {
	server.connections = append(server.connections, client)

	go server.ReceiveFrom(client)
	client.Listen(server.broadcast)
}

func (server *Server) ReceiveFrom(client Client) {
	var msg string
RECEIVE_LOOP:
	for {
		switch websocket.Message.Receive(client.connection, &msg) {
		case nil:
			server.broadcast <- msg
		case io.EOF:
			server.RemoveClient(client)
			break RECEIVE_LOOP
		}
	}
}

// Close connection and remove from pool
func (server *Server) RemoveClient(client Client) {
	for i, conn := range server.connections {
		if conn == client {
			server.connections = append(server.connections[:i], server.connections[i+1:]...)
			client.Terminate()
			break
		}
	}
}
