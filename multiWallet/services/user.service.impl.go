package services

import (
	"WBA/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImplement struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) (UserService, error) {
	return &UserServiceImplement{
		usercollection: usercollection,
		ctx:            ctx,
	}, nil
}

func (u *UserServiceImplement) CheckUser(email string) (*models.User, error) {
	var user *models.User
	filter := bson.M{"email": email}
	if err := u.usercollection.FindOne(u.ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}
