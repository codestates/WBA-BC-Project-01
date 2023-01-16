package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ObjectID primitive.ObjectID `json:"_id" bson:"_id" example:"63a73a1c8d989838729bc114"` //기본키
	Email    string             `json:"email" bson:"email"`                                //SNS ID
	NickName string             `json:"nickname" bson:"age"`                               //별명
	Address  string             `json:"address" bson:"address"`                            //지갑주소
	SnsType  SNS_Type           `json:"snstype" bson:"snstype"`                            //SNS 타입
}

type SNS_Type int

const (
	Google SNS_Type = iota
	Kakao
)
