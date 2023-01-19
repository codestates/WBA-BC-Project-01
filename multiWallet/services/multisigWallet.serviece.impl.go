package services

import (
	"WBA/models"
	"WBA/utils"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/mongo"

	multisig "WBA/contracts/multisigwallet"
)

type MultiSigWalletServiceImplement struct {
	uc *mongo.Collection //유저 컬렉션
	wc *mongo.Collection //키저장 컬렉션
	mc *mongo.Collection //멀티시그 컬렉션

	mod *models.Model

	ctx context.Context
}

func NewMultiSigWalletService(userColl *mongo.Collection, walletColl *mongo.Collection, multiColl *mongo.Collection, ctx context.Context, mod *models.Model) (MultiSigWalletService, error) {
	return &MultiSigWalletServiceImplement{uc: userColl, wc: walletColl, mc: multiColl, ctx: ctx, mod: mod}, nil
}

// 멀티시그 만들기 - 컨트랙트 배포 - 배포된 트랜잭션 주소가 곧 지갑주소
func (m *MultiSigWalletServiceImplement) CreateMultiSigWallet(email string, password string, _owners []string, _numConfirmRequired uint) (string, error) {

	//개인키 가져오기 타입 string
	privateKeyString, err := utils.GetPrivateKey(m.wc, email, password)
	if err != nil {
		log.Fatal(err)
	}
	//개인키 ECDSA 타입으로 변경
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		log.Fatal(err)
	}
	//공개키 추출
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	client := m.mod.WemixClient
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)                // in wei
	auth.GasLimit = uint64(2557770)           // in units
	auth.GasFeeCap = big.NewInt(100000000001) //최대 요금
	auth.GasTipCap = big.NewInt(100000000000) //최대 우선 요금

	fmt.Println(gasPrice)
	//문자열인 지갑 주소 common 타입으로 변경
	var owners []common.Address
	for _, adress := range _owners {
		owners = append(owners, common.HexToAddress(adress))
	}
	owners = append(owners, fromAddress)
	//컴펌 요구 개수 타입 변경
	numConfirmRequired := big.NewInt(int64(_numConfirmRequired))

	address, tx, instance, err := multisig.DeployContracts(auth, client, owners, numConfirmRequired)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	_ = instance
	return address.String(), nil
}
