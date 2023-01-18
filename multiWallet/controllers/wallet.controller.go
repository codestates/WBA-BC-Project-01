package controllers

import (
	"WBA/config"
	"WBA/models"
	"WBA/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	walletService services.WalletService
	mod           *models.Model
}

func NewWalletController(ws services.WalletService, config *config.Config, rep *models.Model) (WalletController, error) {
	return WalletController{walletService: ws, mod: rep}, nil
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

	address, _, email := wc.walletService.NewWallet(&walletRequest)
	ctx.HTML(http.StatusOK, "index.html", gin.H{"address": address, "email": email, "isLogined": true})
}

func (wc *WalletController) BalanceTokens(ctx *gin.Context) {
	fmt.Println("[BalanceTokens]")
	accountAddress := ctx.Query("address")
	coinInfos, tokenInfos := wc.walletService.BalanceTokens(accountAddress)
	ctx.JSON(http.StatusOK, gin.H{"coinInfos :": coinInfos, "tokenInfos :": tokenInfos})
}
