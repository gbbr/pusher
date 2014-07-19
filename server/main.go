package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	connection *websocket.Conn
	doneChan   chan bool
}

type Server struct {
	connections []Client
	comm        chan string
}

// Listen to WebSocket connections and register clients to
// communication channels
func (server *Server) Start() {
	server.connections = make([]Client, 0, 10)
	server.comm = make(chan string)

	onConnected := func(ws *websocket.Conn) {
		defer ws.Close()
		server.RegisterConnection(ws)
	}

	http.Handle("/pipe", websocket.Handler(onConnected))
}

// Receive message from connection and send to communication channel
func (server *Server) RegisterConnection(ws *websocket.Conn) {
	var msg string
	client := Client{ws, make(chan bool)}
	server.connections = append(server.connections, client)

	fmt.Println("Connected.")

	go func() {
	CLIENT_LOOP:
		for {
			select {
			case inc := <-server.comm:
				websocket.Message.Send(ws, inc)
			case <-client.doneChan:
				break CLIENT_LOOP
			}
		}
	}()

	for {
		err := websocket.Message.Receive(ws, &msg)

		if err == nil {
			server.comm <- msg
		}

		if err == io.EOF {
			server.CloseConnection(ws)
			break
		}
	}
}

// Close connection and remove from pool
func (server *Server) CloseConnection(ws *websocket.Conn) {
	for i, client := range server.connections {
		if client.connection == ws {
			client.doneChan <- true
			server.connections = append(server.connections[:i], server.connections[i+1:]...)
			ws.Close()
			break
		}
	}
}

func main() {
	server := new(Server)
	go server.Start()
	http.Handle("/", http.FileServer(http.Dir("client")))

	if err := http.ListenAndServe(":888", nil); err != nil {
		fmt.Println("Error initiating file server. Maybe you lack permissions?")
	}
}
