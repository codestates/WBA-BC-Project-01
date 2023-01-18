package services

import (
	"WBA/contracts"
	"WBA/logger"
	"WBA/models"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/sha3"
)

type WalletServiceImplement struct {
	wc  *mongo.Collection
	uc  *mongo.Collection
	ctx context.Context
	mod *models.Model
}

// 심볼 나중에 수정 예정 화면 출력할 때 토큰의 네트워크를 알기 위해
var WemixSymbol string = "WEMIX"
var EthSymbol string = "ETH"
var KlaySymbol string = "KLAY"

func NewWalletService(walletcollection *mongo.Collection, usercollection *mongo.Collection, ctx context.Context, mod *models.Model) (WalletService, error) {

	return &WalletServiceImplement{
		wc:  walletcollection,
		uc:  usercollection,
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

	var userkeystore models.Keystores
	userkeystore.Email = walletDTO.Email
	json.Unmarshal(keyjson, &userkeystore.KeyStore)
	w.wc.InsertOne(w.ctx, userkeystore)

	var user models.User
	user.Address = address
	user.Email = walletDTO.Email
	w.uc.InsertOne(w.ctx, user)

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

func (w *WalletServiceImplement) TransferTokens(mod models.TransferData) models.TransferData {

	fmt.Println("[service.TransferCoin]")

	// //테스트용 주소
	// mod.Network = "WEMIX"
	// mod.FromAddress = "0x314613c08Cb38e3d782688e86f61a563D8959574"
	// mod.ToAddress = "0x3baac3212c76Cdd3189E3d08D3323C3200B22b9d"
	// mod.TokenContract = "0xF01E78a83F3860433B4ef1b1A7f80B82a269A6fd"
	// mod.SendValue = 5

	fmt.Println("mod.Network : ", mod.Network)
	fmt.Println("mod.FromAddress : ", mod.FromAddress)
	fmt.Println("mod.ToAddress : ", mod.ToAddress)
	fmt.Println("mod.TokenContract : ", mod.TokenContract)
	fmt.Println("mod.SendValue : ", mod.SendValue)

	var client *ethclient.Client
	if mod.Network == "WEMIX" {
		client = w.mod.WemixClient
	} else if mod.Network == "ETH" {
		client = w.mod.EthClient
	} else if mod.Network == "KLAY" {
		client = w.mod.KlaytnClient
	} else {
		panic("Network err")
	}

	//작업중
	//const decimal = 10 ^ 18
	value := int64(mod.SendValue * 1000000000000000000)
	//value := int64(mod.SendValue * decimal)
	fmt.Println("value : ", value)
	//value := big.NewInt(1000000000000000000) // in wei (1 eth)

	fromAddress := common.HexToAddress(mod.FromAddress)
	toAddress := common.HexToAddress(mod.ToAddress)

	sendValue := big.NewInt(value)

	var data []byte
	//토큰의 경우 data에 Set 한다.
	if mod.TokenContract != "" {
		tokenAddress := common.HexToAddress(mod.TokenContract)
		data = SetContractData(toAddress, sendValue)
		sendValue = big.NewInt(0) //토큰 전송의 경우, 컨트랙트에 갯수 설정
		toAddress = tokenAddress  //받는 사람에게 토큰 주소로 변경
	}

	//트랜잭션 실행
	tx := StartTransaction(client, fromAddress, toAddress, sendValue, data)
	mod.TransactionInfo = tx
	return mod
}

// 토큰의 경우, 컨트랙트 설정
func SetContractData(toAddress common.Address, sendValue *big.Int) []byte {
	// 컨트랙트 전송시 사용할 함수명
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	paddedAmount := common.LeftPadBytes(sendValue.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	//컨트랙트 전송 정보 입력
	var pdata []byte
	pdata = append(pdata, methodID...)
	pdata = append(pdata, paddedAddress...)
	pdata = append(pdata, paddedAmount...)
	return pdata
}

func StartTransaction(client *ethclient.Client, fromAddress common.Address, toAddress common.Address, sendValue *big.Int, data []byte) string {

	//프라이베잇 키를 가져와야한다.
	fromPrivateKey, err := crypto.HexToECDSA("-")
	if err != nil {
		fmt.Println(err)
	}

	// 현재 계정의 nonce를 가져옴. 다음 트랜잭션에서 사용할 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println(err)
	}

	// 전송할 양, gasLimit, gasPrice 설정. 추천되는 gasPrice를 가져옴
	gasLimit := uint64(41000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	// 트랜잭션 생성
	tx := types.NewTransaction(nonce, toAddress, sendValue, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	// 트랜잭션 서명
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromPrivateKey)
	if err != nil {
		fmt.Println(err)
	}

	// RLP 인코딩 전 트랜잭션 묶음. 현재는 1개의 트랜잭션
	ts := types.Transactions{signedTx}
	// RLP 인코딩
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	rawTxHex := hex.EncodeToString(rawTxBytes)
	rTxBytes, err := hex.DecodeString(rawTxHex)
	if err != nil {
		fmt.Println(err.Error())
	}

	// RLP 디코딩
	rlp.DecodeBytes(rTxBytes, &tx)
	// 트랜잭션 전송
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println(err)
	}
	//출력된 tx.hash를 익스플로러에 조회 가능
	//예) 0x4788935cfa4a0f23807ba7d7b17a6304cc52795616889fdb9ebdb4498adf4a35
	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())
	return tx.Hash().Hex()
}

func (w *WalletServiceImplement) BalanceTokens(address string) ([]models.TokenInfo, []models.TokenInfo) {

	logger.Info("[WalletServiceImplement.BalanceTokens]")

	var coinInfos []models.TokenInfo
	coinInfos = append(coinInfos, SetCoinInfo(WemixSymbol, address, w.mod.WemixClient))
	coinInfos = append(coinInfos, SetCoinInfo(EthSymbol, address, w.mod.EthClient))
	coinInfos = append(coinInfos, SetCoinInfo(KlaySymbol, address, w.mod.KlaytnClient))

	var tokenInfos []models.TokenInfo
	tokenInfos = SetTokenInfo(WemixSymbol, address, w.mod.WemixClient, w.mod.WemixTokenAddress, tokenInfos)
	tokenInfos = SetTokenInfo(EthSymbol, address, w.mod.EthClient, w.mod.EthTokenAddress, tokenInfos)
	tokenInfos = SetTokenInfo(KlaySymbol, address, w.mod.KlaytnClient, w.mod.KlaytnTokenAddress, tokenInfos)

	return coinInfos, tokenInfos
}

func GetEthValue(balance *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}

func SetCoinInfo(symbolName string, accountAddress string, client *ethclient.Client) models.TokenInfo {
	account := common.HexToAddress(accountAddress)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		logger.Error(err)
	}
	var coinInfo models.TokenInfo
	coinInfo.SymbolName = symbolName
	coinInfo.Network = symbolName
	coinInfo.BalanceOf = GetEthValue(balance)
	return coinInfo
}

func SetTokenInfo(symbolName string, accountAddress string, client *ethclient.Client, contractAddress map[string]string, tokenInfos []models.TokenInfo) []models.TokenInfo {
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

func BalanceToken(client *ethclient.Client, ownerAddress string, contranct string, network string) (models.TokenInfo, error) {

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

	return tokenInfo, err
}
func (w *WalletServiceImplement) TrackAddress(from string) []models.Transaction {
	filter := bson.M{"from": from}
	opts := options.Find().SetSort(bson.D{{"blocknumber", 1}})
	cursor, err := w.wc.Find(context.TODO(), filter, opts)
	var transactions []models.Transaction

	if err = cursor.All(context.TODO(), &transactions); err != nil {
		panic(err)
	}

	for _, result := range transactions {
		cursor.Decode(&result)
		output, err := json.MarshalIndent(result, "", "   ")
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
	}

	return transactions

}
