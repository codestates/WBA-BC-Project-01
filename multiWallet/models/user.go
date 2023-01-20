package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //기본키
	Email    string             `json:"email" bson:"email"`                 //SNS ID
	Address  string             `json:"address" bson:"address"`             //지갑주소
}
