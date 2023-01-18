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

// 자신이 가지고 있는 코인/토큰 정보
type TokenInfo struct {
	Contract   string     `bson:"contract"`
	SymbolName string     `bson:"symbolName"`
	BalanceOf  *big.Float `bson:"balanceOf"`
	Price      int64      `bson:"price"`
	Network    string     `bson:"network"` //코인 이름을 저장 (출력시 참고)
}

// 코인전송 데이터
type TransferData struct {
	Network         string  `bson:"network"`         //
	TokenName       string  `bson:"tokenName"`       //
	FromAddress     string  `bson:"address"`         //
	ToAddress       string  `bson:"address"`         //
	TokenBalance    float64 `bson:"tokenBalance"`    //
	TokenContract   string  `bson:"tokenContract"`   //
	SendValue       float64 `bson:"sendValue"`       //
	PrivateKey      string  `bson:"privateKey"`      //
	TransactionInfo string  `bson:"transactionInfo"` //

}
