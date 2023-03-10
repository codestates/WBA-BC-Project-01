package services

import (
	"WBA/models"
	"crypto/ecdsa"
)

type WalletService interface {
	NewMnemonic() (*models.MnemonicResponse, error)                                                                                                //니모닉 코드 생성
	NewWallet(*models.WalletCreateRequest) (string, *ecdsa.PrivateKey, string)                                                                     //지갑 생성 (개인키 생성)
	NewWalletWithKeystore(privateKey *ecdsa.PrivateKey, address string, walletDTO *models.WalletCreateRequest) (string, *ecdsa.PrivateKey, string) //패스워드로 개인키 키스토어 파일 추출
	BalanceTokens(string) ([]models.TokenInfo, []models.TokenInfo)
	TransferTokens(models.TransferData) models.TransferData
	TrackByAddress(from string) models.TxInfo //FROM 주소 받아 트랜잭션 반환
	TrackByContract(to string) models.TxInfo  //tO 주소 받아 트랜잭션 반환
}
