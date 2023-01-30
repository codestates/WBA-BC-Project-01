package models

import (
	"WBA/ethDaemon/config"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client   *mongo.Client
	colBlock *mongo.Collection
}

func NewModel(cfg *config.Config) (*Model, error) {
	r := &Model{}

	var err error
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.DB.Host)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database(cfg.DB.Database)
		r.colBlock = db.Collection(cfg.DB.Collection)
	}
	return r, nil
}

func (p *Model) SaveBlock(transaction Transaction) error {
	result, err := p.colBlock.InsertOne(context.TODO(), transaction)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Insert succeed: %s\n", result.InsertedID)
	return nil
}

type Block struct {
	BlockHash    string        `bson:"blockHash"`
	BlockNumber  uint64        `bson:"blockNumber"`
	GasLimit     uint64        `bson:"gasLimit"`
	GasUsed      uint64        `bson:"gasUsed"`
	Time         uint64        `bson:"timestamp"`
	Nonce        uint64        `bson:"nonce"`
	Transactions []Transaction `bson:"transactions"`
}

type Transaction struct {
	TxHash      string `bson:"hash"`
	From        string `bson:"from"`
	To          string `bson:"to"` // 컨트랙트의 경우 nil 반환
	Nonce       uint64 `bson:"nonce"`
	GasPrice    uint64 `bson:"gasPrice"`
	GasLimit    uint64 `bson:"gasLimit"`
	Amount      uint64 `bson:"amount"`
	BlockHash   string `bson:"blockHash"`
	BlockNumber uint64 `bson:"blockNumber"`
}

var Address []string

func init() {
	Address = append(Address, "0x50f2Ca639b8F2819F977b73807E0e52e05e6bd70", "0x0924Dfa9cA977e6A956f0399F959c495E54152cb", "0xAbfA9765414d42A3A53B43CDA18Cd879c8E22353")
}
