package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Quote struct
type Quote struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Quote   string             `json:"quote" form:"quote" binding:"required" bson:"quote"`
	Class   string             `json:"class" form:"class" binding:"required" bson:"class"`
	Version int                `json:"version" bson:"__v"`
}
