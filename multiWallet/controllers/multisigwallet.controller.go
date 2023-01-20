package controllers

import (
	"WBA/models"
	"WBA/services"
	"fmt"
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
	walletname := ctx.PostForm("walletname")
	var owners []string
	owners = append(owners, ctx.PostForm("onwer1"))
	owners = append(owners, ctx.PostForm("onwer2"))

	numConfirmRequired, err := strconv.Atoi(ctx.PostForm("confirm"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	wallet, walletname, err := m.ms.CreateMultiSigWallet(email, password, owners, uint(numConfirmRequired), walletname)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"다중서명지갑주소": wallet, "지갑 이름": walletname})
}

func (m *MultisigWalletContrller) SubmitTransaction(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	walletaddress := ctx.PostForm("walletaddress")
	_to := ctx.PostForm("to")
	_value, _ := strconv.Atoi(ctx.PostForm("value"))
	_data := ctx.PostForm("data")

	fmt.Println(email, password, walletaddress, _to, _value)
	tx := m.ms.SubmitTransaction(email, password, walletaddress, _to, _value, _data)
	ctx.JSON(http.StatusCreated, gin.H{"Submit Tx ": tx})
}

func (m *MultisigWalletContrller) ConfirmTransaction(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	walletaddress := ctx.PostForm("walletaddress")
	txIndex, _ := strconv.Atoi(ctx.PostForm("txIndex"))
	tx := m.ms.ConfirmTransaction(email, password, walletaddress, txIndex)

	ctx.JSON(http.StatusCreated, gin.H{"Confirm Tx ": tx})
}

func (m *MultisigWalletContrller) GetTransactionCount(ctx *gin.Context) {
	walletaddress := ctx.Param("wallet")
	cnt := m.ms.GetTransactionCount(walletaddress)
	ctx.JSON(http.StatusOK, gin.H{"트랜잭션 개수": cnt})
}

func (m *MultisigWalletContrller) GetOwners(ctx *gin.Context) {
	walletaddress := ctx.Param("wallet")
	onwers := m.ms.GetOwners(walletaddress)
	ctx.JSON(http.StatusOK, gin.H{"소유자 명단": onwers})
}
