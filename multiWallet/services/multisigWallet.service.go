package services

import "github.com/ethereum/go-ethereum/common"

type MultiSigWalletService interface {
	//이메일,패스워드, 주소배열, 컨펌수 를받아 멀티시그지갑 주소를 리턴합니다.
	CreateMultiSigWallet(email string, password string, addresss []string, NumConfirmRequired uint, walletname string) (string, string, error)
	SubmitTransaction(email string, password string, wallet string, _to string, _value int, _data string) string
	GetTransactionCount(wallet string) string
	ConfirmTransaction(email string, password string, wallet string, _txIndex int) string
	GetOwners(wallet string) []common.Address
}

/*
	O 멀티시그지갑 생성하기 - 컨트랙트 배포 (이메일,패스워드, 주소들, 컨펌수) (트랜잭션 해시, 에러)
	O 트랜잭션 제출하기
	O 제출된 트랜잭션 확인하기
	O 소유자들 반환
	O 현재 총 트랜잭션 수 반환
	제출된 트랜잭션 거절하기
	트랜잭션 실행하기
	특정 트랜잭션 상세정보 가져오기
*/
