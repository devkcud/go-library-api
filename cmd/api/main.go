package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURL string = "mongodb://localhost:27017"

type Book struct {
	ID     int    `bson:"_id" binding:"required"`
	Name   string `bson:"name" binding:"required"`
	Author string `bson:"author" binding:"required"`
}

func main() {
	// Establish a connection to MongoDB
	client := connectMongoDB()

	// Create a database
	db := client.Database("library")
	booksCollection := db.Collection("books")

	// Create a gin router
	router := gin.Default()

	router.POST("/books", func(c *gin.Context) {
		var book Book

		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := booksCollection.InsertOne(context.TODO(), book)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	router.Run(":8080")

	// Disconnect from MongoDB
	defer client.Disconnect(context.TODO())
	log.Println("Disconnected from MongoDB!")
}

func connectMongoDB() *mongo.Client {
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
	return client
}
