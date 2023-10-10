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

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	book.ID = primitive.NewObjectID()

	bc.collection.InsertOne(c.Request.Context(), &book)
}

func (bc *BookCollection) GetSpecificBook(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := bc.collection.FindOne(c.Request.Context(), bson.M{"_id": objID})

	var book models.Book

	if err := result.Decode(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bc *BookCollection) GetBooks(c *gin.Context) {
	cursor, err := bc.collection.Find(c.Request.Context(), bson.D{})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var books []models.Book

	if err = cursor.All(c.Request.Context(), &books); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}
