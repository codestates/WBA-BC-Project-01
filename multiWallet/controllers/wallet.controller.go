package controllers

import (
	"WBA/config"
	"WBA/models"
	"fmt"
	"log"
	"net/http"

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
	ctx.JSON(http.StatusOK, result)
}

func (wc *WalletController) NewWallet(ctx *gin.Context) {
	mnemonic := ctx.PostForm("mnemonic")
	fmt.Println(mnemonic)

	seed, _ := hdwallet.NewSeedFromMnemonic(mnemonic)
	wallet, _ := hdwallet.NewFromSeed(seed)
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, _ := wallet.Derive(path, false)
	privateKey, _ := wallet.PrivateKeyHex(account)
	address := account.Address.Hex()
	var result models.WalletResponse
	result.PrivateKey = privateKey
	result.Address = address
	password := ctx.PostForm("password")

	email := ctx.PostForm("email")
	wc.NewWalletWithKeystore(mnemonic, password, email, ctx)
}
func (wc *WalletController) NewWalletWithKeystore(mnemonic, password string, email string, ctx *gin.Context) {
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

	fmt.Println(string(keyjson))

	/* --- .json 파일 생성하는 부분 */
	// keystoreName := strings.Join([]string{address, "json"}, ".")
	// keystorefile := strings.Join([]string{"./config", keystoreName}, "/")
	// if err := ioutil.WriteFile(keystorefile, keyjson, 0700); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	/* 파일생성대신 DB에 저장하는 로직필요 */
	ctx.IndentedJSON(http.StatusOK, gin.H{"공개키": address, "키저장소": privateKey, "이메일": email})
}
