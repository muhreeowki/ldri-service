package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface{}

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoStore(conUri string) (*MongoStore, error) {
	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conUri))
	if err != nil {
		return nil, err
	}
	// Check if connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database("LDRI")

	return &MongoStore{
		client,
		db,
	}, nil
}

func (s *MongoStore) Close() {
	if err := s.client.Disconnect(context.TODO()); err != nil {
		failOnError(err, "failed to close mongo connection")
	}
}
