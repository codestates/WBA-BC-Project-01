package controllers

import (
	"WBA/config"
	"WBA/models"
	"WBA/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	walletService services.WalletService
}

func NewWalletController(ws services.WalletService, config *config.Config) (WalletController, error) {
	return WalletController{walletService: ws}, nil
}

func (wc *WalletController) NewMnemonic(ctx *gin.Context) {
	mnemonic, err := wc.walletService.NewMnemonic()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
	}
	ctx.JSON(http.StatusOK, mnemonic)
}

func (wc *WalletController) NewWallet(ctx *gin.Context) {
	var walletRequest models.WalletCreateRequest
	walletRequest.Mnemonic = ctx.PostForm("mnemonic")
	walletRequest.Password = ctx.PostForm("password")
	walletRequest.Email = ctx.PostForm("email")

	address, privateKey, email := wc.walletService.NewWallet(&walletRequest)
	ctx.IndentedJSON(http.StatusOK, gin.H{"공개키": address, "키저장소": privateKey, "이메일": email})

}
