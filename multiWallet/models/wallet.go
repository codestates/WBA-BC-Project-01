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
	Email    string             `json:"email" bson:"email"`
	KeyStore encryptedKeyJSONV3 `json:"keystore" bson:"keystore"`
}

type WalletCreateRequestWithPassword struct {
	Mnemonic string `json:"mnemonic" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type encryptedKeyJSONV3 struct {
	Address string     `json:"address"`
	Crypto  CryptoJSON `json:"crypto"`
	Id      string     `json:"id"`
	Version int        `json:"version"`
}

type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}
