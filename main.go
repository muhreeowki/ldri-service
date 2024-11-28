package main

import (
	"log"
	"os"
)

func main() {
	// Connect to MessageQueue Server
	queue, err := NewMessageQueue(os.Getenv("AMQP_URI"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer queue.Close()
	log.Printf("Successfully Connected to RabbitMQ: %v", queue)

	// Connect to Storage Server
	store, err := NewMongoStore(os.Getenv("MONGO_URI"))
	failOnError(err, "Failed to connect to MongoDB")
	defer store.Close()
	log.Printf("Successfully Connected to MongoDB: %v", store)

	server := NewServer(":8080", nil, nil)
	log.Fatal(server.ListenAndServer())
}
