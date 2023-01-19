package route

import (
	"WBA/controllers"
	"WBA/docs"
	"WBA/logger"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swgFiles "github.com/swaggo/files"
	ginSwg "github.com/swaggo/gin-swagger"
)

type Router struct {
	cc *controllers.Controller
	lc *controllers.GoogleLoginController
	wc *controllers.WalletController
	mc *controllers.MultisigWalletContrller
}

/* 주문자, 피주문자 컨트롤러 할당 */
func NewRouter(ctl *controllers.Controller, loginCtl *controllers.GoogleLoginController, walletCtl *controllers.WalletController, multisigCtl *controllers.MultisigWalletContrller) (*Router, error) {
	r := &Router{cc: ctl, lc: loginCtl, wc: walletCtl, mc: multisigCtl}

	return r, nil
}

// cross domain을 위해 사용
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//허용할 header 타입에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		//허용할 method에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// 임의 인증을 위한 함수
func liteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c == nil {
			c.Abort() // 미들웨어에서 사용, 이후 요청 중지
			return
		}
		auth := c.GetHeader("Authorization")

		if auth != "codz" {
			//로직 추가 가능 현재는 Print 로만 처리
			fmt.Println("Authorization failed")
		}
		c.Next()
	}
}

// 세션 관리
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("id")
		if sessionID == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
		}
	}
}

func (p *Router) Idx() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	e := gin.Default()

	e.LoadHTMLGlob("templates/*.html")
	e.Static("/static", "./static/")

	e.Use(logger.GinLogger())
	e.Use(logger.GinRecovery(true))
	e.Use(CORS())

	logger.Info("start server")
	e.GET("/health")
	store := cookie.NewStore([]byte("secret"))
	e.Use(sessions.Sessions("mySession", store))

	e.GET("/swagger/:any", ginSwg.WrapHandler(swgFiles.Handler))
	docs.SwaggerInfo.Host = "localhost:8080" //swagger 정보 등록
	docs.SwaggerInfo.Title = "멀티 체인 지갑"
	docs.SwaggerInfo.Description = "멀티 체인 지갑 API"
	docs.SwaggerInfo.Version = "v01"

	/* 기본 페이지 */
	e.GET("/", p.cc.Index)

	/* 로그인 라우팅 */
	login := e.Group("/auth", liteAuth())
	{
		login.GET("/google/login", p.lc.GoogleLoginHandler)
		login.GET("/google/callback", p.lc.GoogleAuthCallback)
	}

	/* 지갑 라우팅 */
	wallet := e.Group("/wallet", liteAuth())
	{
		wallet.GET("/trackAddress/:from", p.wc.TrackByAddress)
		wallet.GET("/trackContract/:to", p.wc.TrackByContract)
		wallet.POST("/mnemonics", p.wc.NewMnemonic)
		wallet.POST("/", p.wc.NewWallet)
		wallet.GET("/balance", p.wc.BalanceTokens)
		wallet.POST("/transfer", p.wc.TransferTokens)
		// wallet.POST("/keystores", p.wc.NewWalletWithKeystore)
	}
	/* 다중서명지갑 라우팅 */
	multisigwallet := e.Group("/multisigwallet", Authentication())
	{
		multisigwallet.GET("/", p.mc.CreateMultiSigWalletPage)           //지갑생성페이지
		multisigwallet.POST("/", p.mc.CreateMultiSigWallet)              //지갑생성
		multisigwallet.POST("/submit", p.mc.SubmitTransaction)           //Tx 실행
		multisigwallet.POST("/confirm", p.mc.ConfirmTransaction)         //Tx 컨펌
		multisigwallet.GET("/txCount/:wallet", p.mc.GetTransactionCount) //Tx 개수
		multisigwallet.GET("/onwers/:wallet", p.mc.GetOwners)            //소유자들 반환
	}
	return e
}
