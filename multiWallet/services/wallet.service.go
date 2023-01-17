package services

import (
	"WBA/models"
	"crypto/ecdsa"
)

type WalletService interface {
	NewMnemonic() (*models.MnemonicResponse, error)
	NewWallet(*models.WalletCreateRequest) (string, *ecdsa.PrivateKey, string)
	NewWalletWithKeystore(privateKey *ecdsa.PrivateKey, address string, walletDTO *models.WalletCreateRequest) (string, *ecdsa.PrivateKey, string)
	BalanceTokens(string) ([]*models.TokenInfo, []*models.TokenInfo)
}
