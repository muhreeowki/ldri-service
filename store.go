package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	GetUser(email, password string) (*User, error)
	CreateUser(email, password string) (*User, error)
	FetchData() ([]byte, error)
	Close()
}

// User represents a user in the system
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

	coll := db.Collection("users")
	_, err = coll.InsertOne(context.Background(), bson.M{"email": "bob@test.com", "password": "password"})
	if err != nil {
		fmt.Println("failed to insert test user")
	} else {
		fmt.Println("inserted test user")
	}

	return &MongoStore{
		client,
		db,
	}, nil
}

func (s *MongoStore) CreateUser(email, password string) (*User, error) {
	coll := s.db.Collection("users")
	_, err := coll.InsertOne(context.Background(), bson.M{"email": email, "password": password})
	if err != nil {
		return nil, err
	}
	return &User{email, password}, nil
}

func (s *MongoStore) GetUser(email, password string) (*User, error) {
	coll := s.db.Collection("users")

	filter := bson.D{{"email", email}}

	res := new(User)
	err := coll.FindOne(context.Background(), filter).Decode(res)
	if err != nil {
		return nil, err
	}

	if res.Password != password {
		return nil, fmt.Errorf("invalid password")
	}

	return res, nil
}

func (s *MongoStore) FetchData() ([]byte, error) {
	return []byte("hello world"), nil
}

func (s *MongoStore) Close() {
	if err := s.client.Disconnect(context.TODO()); err != nil {
		failOnError(err, "failed to close mongo connection")
	}
}
