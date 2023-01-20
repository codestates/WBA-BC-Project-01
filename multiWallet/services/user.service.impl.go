package services

import (
	"WBA/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImplement struct {
	usercollection        *mongo.Collection
	multiwalletcollection *mongo.Collection
	ctx                   context.Context
}

func NewUserService(usercollection *mongo.Collection, mc *mongo.Collection, ctx context.Context) (UserService, error) {
	return &UserServiceImplement{
		usercollection:        usercollection,
		multiwalletcollection: mc,
		ctx:                   ctx,
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

func (u *UserServiceImplement) GetAddress(email string) (string, error) {
	var user *models.User
	if err := u.usercollection.FindOne(u.ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return "", err
	}
	return user.Address, nil
}

func (u *UserServiceImplement) IsExistMultiWallet(email string) (*models.MultiSigWallet, error) {
	var multiwallet *models.MultiSigWallet
	err := u.multiwalletcollection.FindOne(u.ctx, bson.M{"email": email}).Decode(&multiwallet)
	if err != nil {
		return nil, err
	}
	return multiwallet, nil
}
