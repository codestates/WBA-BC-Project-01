package utils

import (
	"WBA/klyDaemon/models"

	"github.com/klaytn/klaytn/blockchain/types"
)

// Kyaltn 이 지원하는 블록 구조체를 받아 DB에 저장하기 위한 블록 구조체로 변환시켜 리턴합니다.
func BindingBlock(block *types.Block) (models.Block, error) {
	var b models.Block

	b.BlockHash = block.Hash().Hex()
	b.BlockNumber = block.Number().Uint64()
	b.GasUsed = block.GasUsed()
	b.Time = block.Time().Uint64()
	b.Nonce = block.HashNoNonce().Big().Uint64()
	b.Transactions = make([]models.Transaction, 0)

	return b, nil
}

/*
func GetTransactionsFromBlock(txs types.Transactions, b *models.Block, block *types.Block) error {
	if len(txs) > 0 {
		for _, tx := range txs {


			//특정 주소에서 발생한 트랜잭션만 감지
			for _, ads := range models.Address {
				if msg.From().Hex() == ads || tx.To().Hex() == ads { //ads = adress(s)
					// 트랜잭션 구조체 생성
					t := models.Transaction{
						TxHash:      tx.Hash().Hex(),
						To:          "", // 디폴트 값 처리
						From:        msg.From().Hex(),
						Nonce:       tx.Nonce(),
						GasPrice:    tx.GasPrice().Uint64(),
						GasLimit:    tx.Gas(),
						Amount:      tx.Value().Uint64(),
						BlockHash:   block.Hash().Hex(),
						BlockNumber: block.Number().Uint64(),
					}

					if tx.To() != nil {
						t.To = tx.To().Hex()
					}

					b.Transactions = append(b.Transactions, t)
				}
			}
		}
	}
	return nil
}

*/
