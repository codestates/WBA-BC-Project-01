package main

import (
	"WBA/config"
	"WBA/controllers"
	"WBA/logger"
	"WBA/models"
	"WBA/route"
	"WBA/services"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

var (
	server      *gin.Engine
	us          services.UserService
	userc       *mongo.Collection
	mongoClient *mongo.Client
	err         error
	g           errgroup.Group
	cf          *config.Config
	rt          *route.Router
)

// 컨트롤러
var (
	cc controllers.Controller
	lc controllers.LoginController
	uc controllers.UserController
)

func init() {
	var configFlag = flag.String("config", "./config/config.toml", "toml file to use for configuration")
	flag.Parse()
	cf = config.NewConfig(*configFlag)

	/* 로그 초기화 */
	if err := logger.InitLogger(cf); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	logger.Debug("ready server....")

	/* MongoDB Connection */
	if mongoClient, err = models.NewModel(cf.DB.Host); err != nil {
		panic(err)
		/* Router 초기화 */
	} else if lc, err = controllers.NewGoogleLoginController(cf); err != nil {
		panic(err)
	} else if cc, err = controllers.NewController(); err != nil {
		panic(err)
	} else if rt, err = route.NewRouter(&cc, &lc); err != nil {
		panic(fmt.Errorf("router.NewRouter > %v", err))
	}

}

func main() {

	/* Server 설정 */
	mapi := &http.Server{
		Addr:           cf.Server.Port,
		Handler:        rt.Idx(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	g.Go(func() error {
		return mapi.ListenAndServe()
	})

	/* 우아한 종료 */
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := mapi.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		logger.Info("timeout of 1 seconds.")
	}

	logger.Info("Server exiting")

	if err := g.Wait(); err != nil {
		logger.Error(err)
	}

}
