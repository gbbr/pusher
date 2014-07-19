package main

import (
	"./server"
	"fmt"
	"net/http"
)

func main() {
	server := server.New()
	go server.Start()
	http.Handle("/", http.FileServer(http.Dir("client")))

	if err := http.ListenAndServe(":888", nil); err != nil {
		fmt.Println("Error initiating file server. Maybe you lack permissions?")
	}
}
