package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Tags         []string           `bson:"tags" json:"tags"`
	Ingredients  []string           `bson:"ingredients" json:"ingredients"`
	Instructions []string           `bson:"instructions" json:"instructions"`
	PublishedAt  time.Time          `bson:"publishedAt" json:"publishedAt"`
}
