package route

import (
	"WBA/controllers"
	"WBA/docs"
	"WBA/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	swgFiles "github.com/swaggo/files"
	ginSwg "github.com/swaggo/gin-swagger"
)

type Router struct {
	cc *controllers.Controller
	lc *controllers.LoginController
}

/* 주문자, 피주문자 컨트롤러 할당 */
func NewRouter(ctl *controllers.Controller, loginCtl *controllers.LoginController) (*Router, error) {
	r := &Router{cc: ctl, lc: loginCtl}

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

	return e
}
