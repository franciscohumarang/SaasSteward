package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `json:"name,omitempty"`
}
