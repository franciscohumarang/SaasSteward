package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BillingDetail struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	BillingID   primitive.ObjectID `bson:"billing_id,omitempty"`
	Description string             `json:"description,omitempty"`
	Amount      float64            `json:"amount,omitempty"`
}
