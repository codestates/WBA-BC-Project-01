package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	Client   *mongo.Client
	colBlock *mongo.Collection
}

func NewModel(mongoUrl string) (*mongo.Client, error) {
	r := &Model{}

	var err error
	mgUrl := mongoUrl
	if r.Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.Client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.Client.Database("Users")
		r.colBlock = db.Collection("user-info")
	}

	fmt.Println("Mongo DB Successful Connected")

	return r.Client, nil
}
