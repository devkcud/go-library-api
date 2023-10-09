package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURL string = "mongodb://localhost:27017"

func main() {
	// Establish a connection to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Create a collection
	collection := client.Database("test").Collection("users")

	// Insert a document
	result, err := collection.InsertOne(context.TODO(), bson.D{
		{"name", "Pato"},
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted a single document: ", result.InsertedID)

	// Disconnect from MongoDB
	defer client.Disconnect(context.TODO())
	log.Println("Disconnected from MongoDB!")
}
