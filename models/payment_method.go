package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentMethod struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty" json:"id,omitempty"`
	CreditCardNo string             `json:"credit_card_number,omitempty"`
	Expiry       string             `json:"expiry,omitempty"`
	CCV          string             `json:"ccv,omitempty"`
	IsDefault    bool               `json:"is_default,omitempty"`
}
