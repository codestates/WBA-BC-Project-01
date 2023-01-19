package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	conf "WBA/klyDaemon/config"
	logger "WBA/klyDaemon/logger"
	"WBA/klyDaemon/models"
	"WBA/klyDaemon/utils"

	klaytypes "github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
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
	client, err := client.Dial(cf.Klaytn.Url)
	if err != nil {
		log.Fatal(err)
	} else {
		logger.Info("Transaction Daemon Server Start")
		log.Println("Transaction Daemon Server Start")
	}
	headers := make(chan *klaytypes.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				log.Fatal(err)
			}

			/* 블록 구조체 생성 */
			b, err := utils.BindingBlock(block)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Klaytn New BlockNumber : %d ", b.BlockNumber)
			fmt.Printf("%+v", b)
			/* 트랜잭션 추출 */
			//err = utils.GetTransactionsFromBlock(block.Transactions(), &b, block)
			//if err != nil {
			//	log.Fatal(err)
			//}

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
