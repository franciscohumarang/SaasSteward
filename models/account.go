package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	CategoryID    primitive.ObjectID `bson:"category_id,omitempty"`
	Name          string             `json:"name,omitempty"`
	Description   string             `json:"description,omitempty"`
	URL           string             `json:"url,omitempty"`
	Username      string             `json:"username,omitempty"`
	Password      string             `json:"password,omitempty"`
	DateSubscribe string             `json:"date_subscribe,omitempty"`
	BillingType   string             `json:"billing_type,omitempty"`
	BillingDate   string             `json:"billing_date,omitempty"`
	AlertMe       bool               `json:"alert_me,omitempty"`
	IsAutoRenew   bool               `json:"is_auto_renew,omitempty"`
}
