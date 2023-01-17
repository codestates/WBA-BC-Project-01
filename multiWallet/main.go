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
	ws          services.WalletService
	userc       *mongo.Collection
	walletc     *mongo.Collection
	mongoClient *mongo.Client
	err         error
	g           errgroup.Group
	cf          *config.Config
	rt          *route.Router
)

// 컨트롤러
var (
	cc controllers.Controller
	lc controllers.GoogleLoginController
	uc controllers.UserController
	wc controllers.WalletController
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
	} else if userc = mongoClient.Database(cf.DB.Database).Collection("member"); err != nil {
		panic(err)
		/* 서비스 초기화 */
	} else if walletc = mongoClient.Database(cf.DB.Database).Collection("wallet"); err != nil {
		panic(err)
	} else if us, err = services.NewUserService(userc, context.TODO()); err != nil {
		panic(err)
	} else if ws, err = services.NewWalletService(walletc, context.TODO()); err != nil {
		panic(err)
		/* 컨트롤러 초기화 */
	} else if lc, err = controllers.NewGoogleLoginController(us, cf); err != nil {
		panic(err)
	} else if cc, err = controllers.NewController(); err != nil {
		panic(err)
	} else if wc, err = controllers.NewWalletController(ws, cf); err != nil {
		panic(err)
	} else if rt, err = route.NewRouter(&cc, &lc, &wc); err != nil {
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
