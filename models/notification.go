package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty"`
	AccountID                primitive.ObjectID `bson:"account_id,omitempty"`
	DaysToAlertBeforeExpiry  int                `json:"days_to_alert_before_expiry,omitempty"`
	HoursToAlertBeforeExpiry int                `json:"hours_to_alert_before_expiry,omitempty"`
}
