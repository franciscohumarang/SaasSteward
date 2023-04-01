package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Billing struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	DateCreated   string             `json:"date_created,omitempty"`
	AmountDue     float64            `json:"amount_due,omitempty"`
	DueDate       string             `json:"due_date,omitempty"`
	TransactionID string             `json:"transaction_id,omitempty"`
	ModifiedDate  string             `json:"modified_date,omitempty"`
	Status        string             `json:"status,omitempty"`
}
