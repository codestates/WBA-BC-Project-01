package services

import (
	"WBA/contracts"
	"WBA/logger"
	"WBA/models"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"go.mongodb.org/mongo-driver/mongo"
)

type WalletServiceImplement struct {
	wc  *mongo.Collection
	ctx context.Context
	mod *models.Model
}

func NewWalletService(walletcollection *mongo.Collection, ctx context.Context, mod *models.Model) (WalletService, error) {
	return &WalletServiceImplement{
		wc:  walletcollection,
		ctx: ctx,
		mod: mod,
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

func (w *WalletServiceImplement) BalanceTokens(address string) ([]*models.TokenInfo, []*models.TokenInfo) {
	// 심볼 나중에 수정 예정 화면 출력할 때 토큰의 네트워크를 알기 위해
	SymbolWemix := "WEMIX"
	SymbolEth := "ETH"
	//SymbolKlay := "KLAY"
	fmt.Println("[WalletServiceImplement.BalanceTokens]")

	var coinInfos []*models.TokenInfo
	coinInfos = append(coinInfos, w.SetCoinInfo(SymbolWemix, address, w.mod.WemixClient))
	coinInfos = append(coinInfos, w.SetCoinInfo(SymbolEth, address, w.mod.EthClient))
	//coinInfos = append(coinInfos, w.SetCoinInfo(SymbolKlay, address, w.mod.KlaytnClient))

	var tokenInfos []*models.TokenInfo
	tokenInfos = w.SetTokenInfo(SymbolWemix, address, w.mod.WemixClient, w.mod.WemixTokenAddress, tokenInfos)
	tokenInfos = w.SetTokenInfo(SymbolEth, address, w.mod.EthClient, w.mod.EthTokenAddress, tokenInfos)
	//tokenInfos = SetTokenInfo(SymbolKlay, address, w.mod.KlaytnClient, w.mod.KlaytnTokenAddress, tokenInfos)

	return coinInfos, tokenInfos
}

func GetEthValue(balance *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}

func (w *WalletServiceImplement) SetCoinInfo(symbolName string, accountAddress string, client *ethclient.Client) *models.TokenInfo {
	account := common.HexToAddress(accountAddress)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		logger.Error(err)
	}
	var coinInfo models.TokenInfo
	coinInfo.SymbolName = symbolName
	coinInfo.Network = symbolName
	coinInfo.BalanceOf = GetEthValue(balance)
	return &coinInfo
}

func (w *WalletServiceImplement) SetTokenInfo(symbolName string, accountAddress string, client *ethclient.Client, contractAddress map[string]string, tokenInfos []*models.TokenInfo) []*models.TokenInfo {
	for _, contranct := range contractAddress {
		tokenInfo, err := BalanceToken(client, accountAddress, contranct, symbolName)
		if err != nil {
			logger.Error(err)
		} else {
			tokenInfos = append(tokenInfos, tokenInfo)
		}
	}
	return tokenInfos
}

func BalanceToken(client *ethclient.Client, ownerAddress string, contranct string, network string) (*models.TokenInfo, error) {

	tknAddress := common.HexToAddress(contranct)
	ownAddress := common.HexToAddress(ownerAddress)
	instance, err := contracts.NewContracts(tknAddress, client)

	if err != nil {
		logger.Error(err)
	}
	tokenBalance, err := instance.BalanceOf(&bind.CallOpts{}, ownAddress)
	if err != nil {
		logger.Error(err)
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		logger.Error(err)
	}
	var tokenInfo models.TokenInfo
	tokenInfo.SymbolName = symbol
	tokenInfo.Contract = contranct
	tokenInfo.Network = network
	tokenInfo.BalanceOf = GetEthValue(tokenBalance)

	return &tokenInfo, err
}
