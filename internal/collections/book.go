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

func NewBookCollection(db *mongo.Database, collectionName string) *BookCollection {
	return &BookCollection{db.Collection(collectionName)}
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

func (bc *BookCollection) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	// Check if the id is a valid ObjectID
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete the book
	_, err = bc.collection.DeleteOne(c.Request.Context(), bson.M{"_id": objID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a message
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (bc *BookCollection) UpdateBook(c *gin.Context) {
	id := c.Param("id")

	// Check if the id is a valid ObjectID
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bind the book struct with the json request body
	var book models.Book

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the _id field as a the previous ObjectID
	book.ID = objID

	// Update the book
	_, err = bc.collection.UpdateOne(c.Request.Context(), bson.M{"_id": objID}, bson.M{"$set": book})

	c.JSON(http.StatusOK, gin.H{"message": &book})
}
