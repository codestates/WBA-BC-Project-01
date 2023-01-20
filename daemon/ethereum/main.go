package main

import (
	"context"
	"flag"
	"log"

	conf "WBA/ethDaemon/config"
	logger "WBA/ethDaemon/logger"
	"WBA/ethDaemon/models"
	"WBA/ethDaemon/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	cf  *conf.Config
	md  *models.Model
	err error
)

var ()

func init() {
	var configFlag = flag.String("config", "./config/config.toml", "toml file to use for configuration")
	flag.Parse()
	cf = conf.NewConfig(*configFlag)
	/* 로그 설정 */
	logger.InitLogger(cf)

	/* model 초기화 */
	md, err = models.NewModel(cf)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	client, err := ethclient.Dial(cf.Ethereum.Url)

	if err != nil {
		log.Fatal(err)
	} else {
		logger.Info("Transaction Daemon Server Start")
		log.Println("Transaction Daemon Server Start")
	}

	/* 배포한 컨트랙트의 이벤트를 구독할 주소를 입력 */
	contractAddress := common.HexToAddress(cf.ContractAddress.Ca)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	/* 블록 감지 */
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	/* SubscribeFilterLogs 와 쿼리 옵션을 사용하여 이벤트를 감지, 감지한 이벤트트는 인자로 받은 채널에 출력됨 */
	logsContract := make(chan types.Log)
	subContract, err := client.SubscribeFilterLogs(context.Background(), query, logsContract)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-subContract.Err():
			log.Fatal(err)
		case vLog := <-logsContract:
			TransactionLog, _ := utils.Log(vLog).MarshalJSON()
			logger.Event(string(TransactionLog))
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers: // <-- 이더 (goeril)
			block, err := client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				log.Fatal(err)
			}
			/* 블록 구조체 생성 */
			b, err := utils.BindingBlock(block)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Ether New BlockNumber : %d ", b.BlockNumber)
			/* 트랜잭션 추출 */
			err = utils.GetTransactionsFromBlock(block.Transactions(), &b, block)
			if err != nil {
				log.Fatal(err)
			}

			/* 트랜잭션이 존재할 경우만 DB에 저장 */
			if len(b.Transactions) > 0 {
				for _, tr := range b.Transactions {
					if err := md.SaveBlock(tr); err != nil {
						log.Fatal(err)
					}
				}

			}
		}
	}

}
