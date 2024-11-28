package main

import (
	"log"
	"os"
)

func main() {
	queue, err := NewMessageQueue(os.Getenv("AMQP_URI"))
	failOnError(err, "Failed to connect to RabbitMQ")

	log.Printf("Successfully Connected to RabbitMQ: %+v", queue)

	server := NewServer(":8080", nil, nil)
	log.Fatal(server.ListenAndServer())
}
