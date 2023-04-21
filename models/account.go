package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	CompanyID     primitive.ObjectID `bson:"company_id,omitempty" json:"id,omitempty"`
	CategoryID    primitive.ObjectID `bson:"category_id,omitempty"`
	CategoryName  string             `json:"category_name,omitempty"`
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

type AccoutByCategory struct {
	CategoryID   primitive.ObjectID `bson:"category_id,omitempty"`
	CategoryName string             `json:"category_name,omitempty"`
	Accounts     []Account          `json:"accounts,omitempty"`
}
