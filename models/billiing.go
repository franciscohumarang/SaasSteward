package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Billing struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        primitive.ObjectID `bson:"user_id,omitempty" json:"id,omitempty"`
	DateCreated   time.Time          `json:"date_created,omitempty"`
	Description   string             `json:"description,omitempty"`
	AmountDue     float64            `json:"amount_due,omitempty"`
	DueDate       string             `json:"due_date,omitempty"`
	TransactionID string             `json:"transaction_id,omitempty"`
	ModifiedDate  string             `json:"modified_date,omitempty"`
	Status        string             `json:"status,omitempty"`
}
