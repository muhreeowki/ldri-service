package main

import (
	"log"
)

func main() {
	// Connect to MessageQueue Server
	// queue, err := NewMessageQueue(os.Getenv("AMQP_URI"))
	// failOnError(err, "Failed to connect to RabbitMQ")
	// defer queue.Close()
	// log.Printf("Successfully Connected to RabbitMQ: %v", queue)

	// Connect to Storage Server
	// store, err := NewMongoStore(os.Getenv("MONGO_URI"))
	// failOnError(err, "Failed to connect to MongoDB")
	// defer store.Close()
	// log.Printf("Successfully Connected to MongoDB: %v", store)

	store, err := NewMongoStore("mongodb://localhost:27017")
	failOnError(err, "Failed to connect to Mongo")

	// Start the authSvr
	authSvr := NewAuthServer(":8080", store)
	dataSvr := NewDataServer(":8000", store)

	go func() { log.Fatal(dataSvr.Run()) }()
	log.Fatal(authSvr.ListenAndServer())
}
