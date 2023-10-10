package main

import (
	"context"
	"fmt"
	"log"

	"github.com/devkcud/go-library-api/internal/collections"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURL string = "mongodb://localhost:27017"
const port int = 8080

func main() {
	// Establish a connection to MongoDB
	client := connectMongoDB()

	// Create a database
	db := client.Database("library")
	booksCollection := collections.NewBookCollection(db)

	// Create a gin router
	router := gin.Default()

	// Routes
	router.GET("/books", booksCollection.GetBooks)
	router.GET("/books/:id", booksCollection.GetSpecificBook)
	router.POST("/books", booksCollection.PostBook)
	router.DELETE("/books/:id", booksCollection.DeleteBook)
    router.PATCH("/books/:id", booksCollection.UpdateBook)

	// Listen and serve on localhost:8080
	router.Run(fmt.Sprintf(":%d", port))

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
