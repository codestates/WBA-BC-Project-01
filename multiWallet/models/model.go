package models

import (
	"WBA/config"
	"WBA/logger"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	Client   *mongo.Client
	colBlock *mongo.Collection

	WemixClient  *ethclient.Client
	EthClient    *ethclient.Client
	KlaytnClient *ethclient.Client

	WemixTokenAddress  map[string]string
	EthTokenAddress    map[string]string
	KlaytnTokenAddress map[string]string
}

func NewModel(cf *config.Config) (*Model, error) {
	r := &Model{}

	var err error
	mgUrl := cf.DB.Host
	if r.Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.Client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.Client.Database("Daemon")
		r.colBlock = db.Collection("wemixBlock")
	}

	fmt.Println("Mongo DB Successful Connected")

	//Wallet ethclient  추가
	r.WemixClient, err = ethclient.Dial(cf.Client.UrlWemix)
	if err != nil {
		logger.Error(err)
	}
	r.EthClient, err = ethclient.Dial(cf.Client.UrlEth)
	if err != nil {
		logger.Error(err)
	}
	r.KlaytnClient, err = ethclient.Dial(cf.Client.UrlKlay)
	if err != nil {
		logger.Error(err)
	}

	r.WemixTokenAddress = cf.Client.WemixTokenAddress
	r.EthTokenAddress = cf.Client.EthTokenAddress
	r.KlaytnTokenAddress = cf.Client.KlaytnTokenAddress

	fmt.Println("ethclient Successful Connected")

	return r, nil
}
