package controllers

import (
	"WBA/models"
	"WBA/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type MultisigWalletContrller struct {
	us  services.UserService
	ms  services.MultiSigWalletService
	ws  services.WalletService
	mod *models.Model
}

func NewMultiSigWalletController(us services.UserService, ms services.MultiSigWalletService, ws services.WalletService, mod *models.Model) (MultisigWalletContrller, error) {
	return MultisigWalletContrller{us: us, ms: ms, ws: ws, mod: mod}, nil
}

func (m *MultisigWalletContrller) CreateMultiSigWalletPage(ctx *gin.Context) {
	session := sessions.Default(ctx)
	address, err := m.us.GetAddress(session.Get("email").(string))
	if err != nil {
		log.Fatal(err)
	}
	coinInfo, tokenInfo := m.ws.BalanceTokens(address)
	ctx.HTML(http.StatusOK, "multisigwallet.html", gin.H{"email": session.Get("email"), "coinInfo": coinInfo, "tokenInfo": tokenInfo})
}

func (m *MultisigWalletContrller) CreateMultiSigWallet(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	var owners []string
	owners = append(owners, ctx.PostForm("owners1"), ctx.PostForm("owners2"))

	numConfirmRequired, err := strconv.Atoi(ctx.PostForm("confirm"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	wallet, err := m.ms.CreateMultiSigWallet(email, password, owners, uint(numConfirmRequired))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"다중서명지갑주소": wallet})
}
