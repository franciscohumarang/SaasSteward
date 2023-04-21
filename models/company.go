package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `json:"status,omitempty"`
}
