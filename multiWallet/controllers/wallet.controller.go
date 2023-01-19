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

func (wc *WalletController) TransferTokens(ctx *gin.Context) {
	fmt.Println("[TransferTokens]")
	var params models.TransferData

	if err := ctx.ShouldBind(&params); err == nil {
		fmt.Println("[TransferTokens]params:", params)
		params = wc.walletService.TransferTokens(params)
		if params.TransactionInfo != "" {
			ctx.JSON(http.StatusOK, gin.H{"tx sent :": params.TransactionInfo})
		} else {
			ctx.JSON(http.StatusBadRequest, "Network, FromAddress, ToAddress, SendValue, UserMail, UserPWD 또는 TokenContract를 다시 입력해주세요")
		}
	}
}

func (wc *WalletController) TrackByAddress(ctx *gin.Context) {

	// UI에 따라 string 어떻게 가져올지 구현방법 변경 필요
	//r, _ := models.NewModel(cf.DB.Host)
	from := ctx.Param("from")
	scope := wc.walletService.TrackByAddress(from)

	// for _, scp := range scope {
	// 	txHash := scp.TxHash
	// 	blockNum := scp.BlockNumber
	// 	addressFrom := scp.From
	// 	addressTo := scp.To
	// 	amount := scp.Amount
	// }

	ctx.JSON(http.StatusOK, scope)

}

func (wc WalletController) TrackByContract(ctx *gin.Context) {

	to := ctx.PostForm("to")
	scope := wc.walletService.TrackByContract(to)

	ctx.JSON(http.StatusOK, scope)
}
