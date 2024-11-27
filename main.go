package main

import "log"

func main() {
	server := NewServer(":8080", nil, nil)
	log.Fatal(server.ListenAndServer())
}
