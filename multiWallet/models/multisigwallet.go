package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MultiSigWallet struct {
	ObjectID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //기본키
	Email      string             `json:"email" bson:"email"`                 //소유자
	WalletName string             `json:"walletname" bson:"walletname"`       //지갑주소별명
	Address    string             `json:"address" bson:"address"`             //다중서명지갑주소}
}

type Tx struct {
	To               common.Address
	Value            *big.Int
	Data             []byte
	Executed         bool
	NumConfirmations *big.Int
}
