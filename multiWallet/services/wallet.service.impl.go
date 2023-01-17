package services

import (
	"WBA/models"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WalletServiceImplement struct {
	wc  *mongo.Collection
	ctx context.Context
}

func NewWalletService(walletcollection *mongo.Collection, ctx context.Context) (WalletService, error) {
	return &WalletServiceImplement{
		wc:  walletcollection,
		ctx: ctx,
	}, nil
}

func (w *WalletServiceImplement) NewMnemonic() (*models.MnemonicResponse, error) {
	if entropy, err := hdwallet.NewEntropy(256); err != nil {
		return nil, err
	} else if mnemonic, err := hdwallet.NewMnemonicFromEntropy(entropy); err != nil {
		return nil, err
	} else {
		var result models.MnemonicResponse
		result.Mnemonic = mnemonic
		return &result, nil
	}
}

func (w *WalletServiceImplement) NewWallet(walletDTO *models.WalletCreateRequest) (string, *ecdsa.PrivateKey, string) {
	seed, _ := hdwallet.NewSeedFromMnemonic(walletDTO.Mnemonic)
	wallet, _ := hdwallet.NewFromSeed(seed)
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, _ := wallet.Derive(path, false)
	privateKey, _ := wallet.PrivateKey(account)
	address := account.Address.Hex()

	return w.NewWalletWithKeystore(privateKey, address, walletDTO)
}

func (w *WalletServiceImplement) NewWalletWithKeystore(privateKey *ecdsa.PrivateKey, address string, walletDTO *models.WalletCreateRequest) (string, *ecdsa.PrivateKey, string) {

	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}

	ks := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	keyjson, err := keystore.EncryptKey(ks, walletDTO.Password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var userkeystore models.Keystores
	userkeystore.Email = walletDTO.Email
	json.Unmarshal(keyjson, &userkeystore.KeyStore)
	w.wc.InsertOne(w.ctx, userkeystore)

	return address, privateKey, walletDTO.Email
}

// 이메일과 패스워드 받아서 개인키를 반환합니다.
func (w *WalletServiceImplement) GetPrivateKey(email string, password string) (string, error) {
	var keyjson *models.Keystores
	filter := bson.M{"email": email}
	if err := w.wc.FindOne(w.ctx, filter).Decode(&keyjson); err != nil {
		fmt.Println(err)
		return "", err
	}

	keyjsonToByte, err := json.Marshal(keyjson.KeyStore)
	if err != nil {
		return "", err
	}
	key, err := keystore.DecryptKey(keyjsonToByte, string(password))
	if err != nil {
		return "", err
	}
	//개인키 : hex.EncodeToString(key.PrivateKey.D.Bytes()))
	return hex.EncodeToString(key.PrivateKey.D.Bytes()), nil
}
