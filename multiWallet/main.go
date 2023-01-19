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
	ms          services.MultiSigWalletService
	userc       *mongo.Collection
	walletc     *mongo.Collection
	wemixc      *mongo.Collection
	klaytnc     *mongo.Collection
	ethereumc   *mongo.Collection
	multisigcol *mongo.Collection
	mongoClient *mongo.Client
	err         error
	g           errgroup.Group
	cf          *config.Config
	rt          *route.Router
	mod         *models.Model
)

// 컨트롤러
var (
	cc controllers.Controller
	lc controllers.GoogleLoginController
	uc controllers.UserController
	wc controllers.WalletController
	mc controllers.MultisigWalletContrller
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
	if mod, err = models.NewModel(cf); err != nil {
		panic(err)
	} else if userc = mod.Client.Database(cf.DB.MultiWalletDatabase).Collection(cf.DB.UserInfoColl); err != nil {
		panic(err)
	} else if walletc = mod.Client.Database(cf.DB.MultiWalletDatabase).Collection(cf.DB.WalletInfoColl); err != nil {
		panic(err)
	} else if wemixc = mod.Client.Database(cf.DB.DaemonDatabase).Collection(cf.DB.WemixColl); err != nil {
		panic(err)
	} else if klaytnc = mod.Client.Database(cf.DB.DaemonDatabase).Collection(cf.DB.KlaytnColl); err != nil {
		panic(err)
	} else if ethereumc = mod.Client.Database(cf.DB.DaemonDatabase).Collection(cf.DB.EthColl); err != nil {
		panic(err)
	} else if multisigcol = mod.Client.Database(cf.DB.MultiWalletDatabase).Collection(cf.DB.MultiSigWalletCollection); err != nil {
		panic(err)
		/* 서비스 초기화 */
	} else if us, err = services.NewUserService(userc, multisigcol, context.TODO()); err != nil {
		panic(err)
	} else if ws, err = services.NewWalletService(walletc, userc, wemixc, klaytnc, ethereumc, context.TODO(), mod); err != nil {
		panic(err)
	} else if ms, err = services.NewMultiSigWalletService(userc, walletc, multisigcol, context.TODO(), mod); err != nil {
		panic(err)
		/* 컨트롤러 초기화 */
	} else if lc, err = controllers.NewGoogleLoginController(us, cf); err != nil {
		panic(err)
	} else if cc, err = controllers.NewController(); err != nil {
		panic(err)
	} else if wc, err = controllers.NewWalletController(ws, cf, mod); err != nil {
		panic(err)
	} else if mc, err = controllers.NewMultiSigWalletController(us, ms, ws, mod); err != nil {
		panic(err)
	} else if rt, err = route.NewRouter(&cc, &lc, &wc, &mc); err != nil {
		panic(fmt.Errorf("router.NewRouter > %v", err))
	}
	mongoClient = mod.Client //어디서 쓰실까봐 살려둠
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
