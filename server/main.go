package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	connections []*websocket.Conn
}

// Listen to WebSocket connections and register clients to
// communication channels
func (server *Server) Start() {
	communication := make(chan string)
	server.connections = make([]*websocket.Conn, 0, 10)

	go server.Listen(communication)

	onConnected := func(ws *websocket.Conn) {
		defer ws.Close()
		server.RegisterConnection(ws, communication)
	}

	http.Handle("/pipe", websocket.Handler(onConnected))
}

// Receive message from connection and send to communication channel
func (server *Server) RegisterConnection(ws *websocket.Conn, c chan string) {
	var msg string
	server.connections = append(server.connections, ws)

	fmt.Println("Connected.")

	for {
		err := websocket.Message.Receive(ws, &msg)
		if err == nil && err != io.EOF {
			c <- msg
		}
	}
}

// Wait for incoming messages and broadcast to all connections
func (server *Server) Listen(c chan string) {
	for {
		msg := <-c
		server.Broadcast(msg)
	}
}

// Broadcast message to all connections
func (s *Server) Broadcast(msg string) {
	for _, conn := range s.connections {
		websocket.Message.Send(conn, msg)
	}
}

func main() {
	server := new(Server)
	go server.Start()

	http.Handle("/", http.FileServer(http.Dir("client")))
	err := http.ListenAndServe(":888", nil)

	if err != nil {
		fmt.Println("Error initiating file server")
	}
}
