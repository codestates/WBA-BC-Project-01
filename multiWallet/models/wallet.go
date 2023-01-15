package models

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"

	"github.com/google/uuid"
)

type MnemonicResponse struct {
	Mnemonic string `json:"mnemonic"`
}

type WalletCreateRequest struct {
	Mnemonic string `json:"mnemonic" binding:"required"`
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

type WalletCreateRequestWithPassword struct {
	Mnemonic string `json:"mnemonic" binding:"required"`
	Password string `json:"password" binding:"required"`
}
