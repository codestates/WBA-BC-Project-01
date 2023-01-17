package services

import (
	"WBA/models"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
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

	var keystore models.Keystores
	keystore.Key = string(keyjson)
	keystore.Email = walletDTO.Email
	w.wc.InsertOne(w.ctx, keystore)
	return address, privateKey, walletDTO.Email
}
