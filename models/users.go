package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName    string             `bson:"first_name,omitempty" json:"first_name,omitempty" validate:"required"`
	LastName     string             `bson:"last_name,omitempty" json:"last_name,omitempty" validate:"required"`
	Email        string             `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	Company      string             `bson:"company,omitempty" json:"company,omitempty"`
	CountryCode  string             `bson:"country_code,omitempty" json:"country_code,omitempty" validate:"required"`
	Password     string             `bson:"password,omitempty" json:"password,omitempty" validate:"required,min=8"`
	MobileNumber string             `bson:"mobile_number,omitempty" json:"mobile_number,omitempty" validate:"required"`
}
