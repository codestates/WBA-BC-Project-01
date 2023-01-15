package services

import (
	"context"

	"WBA/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImplement struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImplement{
		usercollection: usercollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImplement) CreateUser(user *models.User) error {
	_, err := u.usercollection.InsertOne(u.ctx, user)
	return err
}

func (u *UserServiceImplement) GetUser(id string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "id", Value: id}}
	err := u.usercollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}
