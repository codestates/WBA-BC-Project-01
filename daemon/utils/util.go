package utils

import (
	"WBA/daemon/models"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex"`
	// hash of the block in which the transaction was included
	BlockHash common.Hash `json:"blockHash"`
	// index of the log in the block
	Index uint `json:"logIndex"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`
}

// MarshalJSON marshals as JSON.
func (l Log) MarshalJSON() ([]byte, error) {
	type Log struct {
		Address     common.Address `json:"address" gencodec:"required"`
		Topics      []common.Hash  `json:"topics" gencodec:"required"`
		Data        hexutil.Bytes  `json:"data" gencodec:"required"`
		BlockNumber uint64         `json:"blockNumber"`
		TxHash      common.Hash    `json:"transactionHash" gencodec:"required"`
		TxIndex     uint           `json:"transactionIndex"`
		BlockHash   common.Hash    `json:"blockHash"`
		Index       uint           `json:"logIndex"`
		Removed     bool           `json:"removed"`
	}
	var enc Log
	enc.Address = l.Address
	enc.Topics = l.Topics
	enc.Data = l.Data
	enc.BlockNumber = l.BlockNumber
	enc.TxHash = l.TxHash
	enc.TxIndex = l.TxIndex
	enc.BlockHash = l.BlockHash
	enc.Index = l.Index
	enc.Removed = l.Removed
	return json.Marshal(&enc)
}

// Go-Ethereum 타입의 블록 구조체를 받아 DB에 저장하기 위한 블록 구조체로 변환시켜 리턴합니다.
func BindingBlock(block *types.Block) (models.Block, error) {
	var b models.Block

	b.BlockHash = block.Hash().Hex()
	b.BlockNumber = block.Number().Uint64()
	b.GasLimit = block.GasLimit()
	b.GasUsed = block.GasUsed()
	b.Time = block.Time()
	b.Nonce = block.Nonce()
	b.Transactions = make([]models.Transaction, 0)
	return b, nil
}

func GetTransactionsFromBlock(txs types.Transactions, b *models.Block, block *types.Block) error {
	if len(txs) > 0 {
		for _, tx := range txs {
			msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), block.BaseFee())
			if err != nil {
				return err
			}

			/* 특정 주소에서 발생한 트랜잭션만 감지 */
			for _, ads := range models.Address {
				if msg.From().Hex() == ads { //ads = adress(s)
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
