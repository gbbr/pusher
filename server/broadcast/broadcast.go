package broadcast

import (
	"code.google.com/p/go.net/websocket"
	"io"
	"net/http"
)

// Initializes a new server
func New() *Server {
	return &Server{make([]*websocket.Conn, 0, 10)}
}

type Server struct {
	connections []*websocket.Conn
}

// Start server on give URL
func (server *Server) Start(path string) {
	onConnected := func(ws *websocket.Conn) {
		defer ws.Close()
		server.connections = append(server.connections, ws)
		server.Register(ws)
	}

	http.Handle(path, websocket.Handler(onConnected))
}

// Receives messages from client and broadcast
func (server *Server) Register(client *websocket.Conn) {
	var msg string
RECEIVE_LOOP:
	for {
		switch websocket.Message.Receive(client, &msg) {
		case nil:
			server.Broadcast(msg)
		case io.EOF:
			server.RemoveClient(client)
			break RECEIVE_LOOP
		}
	}
}

// Broadcast a message to all connections
func (server *Server) Broadcast(msg string) {
	for _, socket := range server.connections {
		go func(ws *websocket.Conn) {
			websocket.Message.Send(ws, msg)
		}(socket)
	}
}

// Close connection and remove from pool
func (server *Server) RemoveClient(socket *websocket.Conn) {
	for i, client := range server.connections {
		if client == socket {
			server.connections = append(server.connections[:i], server.connections[i+1:]...)
			socket.Close()
			break
		}
	}
}
