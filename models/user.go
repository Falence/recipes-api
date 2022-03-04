package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
	Username string `bson:"username,omitempty" json:"username,omitempty"`
}
