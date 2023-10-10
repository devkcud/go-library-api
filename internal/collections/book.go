package collections

import (
	"net/http"

	"github.com/devkcud/go-library-api/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookCollection struct {
	collection *mongo.Collection
}

func NewBookCollection(db *mongo.Database) *BookCollection {
	return &BookCollection{db.Collection("books")}
}

func (bc *BookCollection) PostBook(c *gin.Context) {
	book := models.Book{}

	// Bind the book struct with the json request body
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Set the _id field as a new ObjectID
	book.ID = primitive.NewObjectID()

	// Insert the book in the collection
	bc.collection.InsertOne(c.Request.Context(), &book)
}

func (bc *BookCollection) GetSpecificBook(c *gin.Context) {
	id := c.Param("id")

	// Check if the id is a valid ObjectID
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the book
	result := bc.collection.FindOne(c.Request.Context(), bson.M{"_id": objID})

	var book models.Book

	// Try decode the book
	if err := result.Decode(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the book
	c.JSON(http.StatusOK, book)
}

func (bc *BookCollection) GetBooks(c *gin.Context) {
	// Find all books
	cursor, err := bc.collection.Find(c.Request.Context(), bson.D{})

	// Check if there is an error in the collection
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var books []models.Book

	// Decode the books and set them in the books slice
	if err = cursor.All(c.Request.Context(), &books); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the books
	c.JSON(http.StatusOK, books)
}
