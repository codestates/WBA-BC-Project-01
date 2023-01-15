package controllers

import (
	"WBA/config"
	"WBA/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type WalletController struct {
	md *models.Model
}

func NewWalletController(config *config.Config) (WalletController, error) {
	return WalletController{}, nil
}

func (wc *WalletController) NewMnemonic(ctx *gin.Context) {
	entropy, _ := hdwallet.NewEntropy(256)
	mnemonic, _ := hdwallet.NewMnemonicFromEntropy(entropy)

	var result models.MnemonicResponse
	result.Mnemonic = mnemonic
	ctx.IndentedJSON(http.StatusOK, result)
}

// ./controller/controller.go
func (wc *WalletController) NewWallet(ctx *gin.Context) {
	var body models.WalletCreateRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mnemonic := body.Mnemonic

	seed, _ := hdwallet.NewSeedFromMnemonic(mnemonic)
	wallet, _ := hdwallet.NewFromSeed(seed)
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")

	account, _ := wallet.Derive(path, false)
	privateKey, _ := wallet.PrivateKeyHex(account)

	address := account.Address.Hex()

	var result models.WalletResponse
	result.PrivateKey = privateKey
	result.Address = address

	ctx.IndentedJSON(http.StatusOK, result)
}

func (wc *WalletController) NewWalletWithKeystore(ctx *gin.Context) {
	var body models.WalletCreateRequestWithPassword
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mnemonic := body.Mnemonic
	password := body.Password

	seed, _ := hdwallet.NewSeedFromMnemonic(mnemonic)
	wallet, _ := hdwallet.NewFromSeed(seed)
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")

	account, _ := wallet.Derive(path, false)
	privateKey, _ := wallet.PrivateKey(account)

	address := account.Address.Hex()

	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}

	ks := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	keyjson, err := keystore.EncryptKey(ks, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		log.Fatalf(err.Error())
	}

	keystoreName := strings.Join([]string{address, "json"}, ".")
	keystorefile := strings.Join([]string{"./config", keystoreName}, "/")
	//wc.md.SaveUserInfo(keyjson)

	if err := ioutil.WriteFile(keystorefile, keyjson, 0700); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"result": "ok"})
}
