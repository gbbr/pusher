package main

import (
	"./broadcast"
	"fmt"
	"net/http"
)

func main() {
	server := broadcast.New()
	server.Start("/pipe")

	http.Handle("/", http.FileServer(http.Dir("client")))

	if err := http.ListenAndServe(":888", nil); err != nil {
		fmt.Println("Error initiating file server. Maybe you lack permissions?")
	}
}
