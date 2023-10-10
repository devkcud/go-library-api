package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Name   string             `bson:"name" json:"name" binding:"required"`
	Author string             `bson:"author" json:"author" binding:"required"`
}
