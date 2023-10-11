package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/devkcud/go-library-api/internal/collections"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI string `json:"mongoURI"`

	DatabaseName   string `json:"databaseName"`
	CollectionName string `json:"collectionName"`

	APIHost string `json:"APIHost"`
	APIPort int    `json:"APIPort"`
}

var config Config

func main() {
	// Get the absolute path of root during build phase
	_, file, _, _ := runtime.Caller(0)
	jsonFile, err := os.Open(filepath.Join(filepath.Dir(file), "..", "..", "config.json"))

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	// Read the config
	jsonByteContent, _ := io.ReadAll(jsonFile)

	// Unmarshal the config
	json.Unmarshal(jsonByteContent, &config)

	// Establish a connection to MongoDB
	clientOptions := options.Client().ApplyURI(config.MongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Create a database
	db := client.Database(config.DatabaseName)
	booksCollection := collections.NewBookCollection(db, config.CollectionName)

	// Create a gin router
	router := gin.Default()

	// Routes
	router.GET("/books", booksCollection.GetBooks)
	router.GET("/books/:id", booksCollection.GetSpecificBook)
	router.POST("/books", booksCollection.PostBook)
	router.DELETE("/books/:id", booksCollection.DeleteBook)
	router.PATCH("/books/:id", booksCollection.UpdateBook)

	// Listen and serve
	router.Run(fmt.Sprintf("%s:%d", config.APIHost, config.APIPort))

	// Disconnect from MongoDB
	defer client.Disconnect(context.TODO())
	log.Println("Disconnected from MongoDB!")
}
