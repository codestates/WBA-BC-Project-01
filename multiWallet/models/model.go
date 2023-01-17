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
		db := r.Client.Database("Users")
		r.colBlock = db.Collection("user-info")
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

	// r.WemixTokenAddress = cf.Client.WemixTokenAddress
	// r.EthTokenAddress = cf.Client.EthTokenAddress
	// r.KlaytnTokenAddress = cf.Client.KlaytnTokenAddress
	r.WemixTokenAddress = map[string]string{"WAL": "0x9Fa7F4E848Df29B3F653c47cC12b4c9bBCf2b99c", "LWS": "0xF01E78a83F3860433B4ef1b1A7f80B82a269A6fd"}
	r.EthTokenAddress = map[string]string{"WAL": "0x9Fa7F4E848Df29B3F653c47cC12b4c9bBCf2b99c", "LWS": "0xF01E78a83F3860433B4ef1b1A7f80B82a269A6fd"}
	r.KlaytnTokenAddress = map[string]string{"WAL": "0x9Fa7F4E848Df29B3F653c47cC12b4c9bBCf2b99c", "LWS": "0xF01E78a83F3860433B4ef1b1A7f80B82a269A6fd"}

	fmt.Println("ethclient Successful Connected")

	return r, nil
}
