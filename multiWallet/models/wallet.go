package models

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/google/uuid"
)

type MnemonicResponse struct {
	Mnemonic string `json:"mnemonic"`
}

type WalletCreateRequest struct {
	Mnemonic string `json:"mnemonic" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"emial" binding:"required"`
}

type WalletResponse struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
}

type Key struct {
	Id         uuid.UUID
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
}

type Keystores struct {
	Email string `json:"email" bson:"email"`
	Key   string `json:"key" bson:"key"`
}

type WalletCreateRequestWithPassword struct {
	Mnemonic string `json:"mnemonic" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenInfo struct {
	Contract   string     `bson:"contract"`
	SymbolName string     `bson:"symbolName"`
	BalanceOf  *big.Float `bson:"balanceOf"`
	Price      int64      `bson:"price"`
	Network    string     `bson:"network"` //코인 이름을 저장 (출력시 참고)
}
