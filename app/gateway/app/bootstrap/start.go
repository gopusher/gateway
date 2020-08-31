package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gopusher/gateway/app/gateway/app/api"
	"github.com/gopusher/gateway/app/gateway/app/cfg"
	_ "github.com/gopusher/gateway/app/gateway/app/includes"
	"github.com/gopusher/gateway/app/gateway/app/protocols"
	"github.com/gopusher/gateway/pkg/config"
	"github.com/gopusher/gateway/pkg/helper"
	"github.com/gopusher/gateway/pkg/log"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func Start(cfgFile string, flagSet *pflag.FlagSet) {
	//load cfg
	viper, err := config.LoadConfig(cfg.Config, cfgFile, flagSet, nil)
	if err != nil {
		panic(err)
	}

	//dingtalks := dingtalk.InitConnections(cfg.Config.DingTalk)

	//init logger
	cfg.Config.InitLoggerConfig()
	log.New(cfg.Config.LoggerConfig)
	defer log.Sync()

	log.Info("using cfg file: " + viper.ConfigFileUsed())
	log.Debug("cfg data", zap.String("config_data", helper.ToJsonString(cfg.Config)))

	//redisConnections := redis.InitConnections(cfg.Config.Redis)

	server, _ := protocols.Load(cfg.Config.Protocol())
	defer func() {
		if err := server.LeaveCluster(); err != nil {
			log.Panic("servers.UnRegisterCluster err", zap.Error(err))
		}
	}()
	if err := server.Run(); err != nil {
		log.Panic("servers.Run err", zap.Error(err))
	}
	if err := server.JoinCluster(); err != nil {
		log.Panic("servers.RegisterCluster err", zap.Error(err))
	}

	go api.InitRpcServer(cfg.Config.Node, server, cfg.Config.ApiServer)

	fmt.Println("gopusher gateway finished bootstrap")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	//register for interupt (Ctrl+C) and SIGTERM (docker)
	//todo 增加平滑重启信号相关处理逻辑
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-quit
	log.Info("get signal, start shutdown server ...", zap.Any("signal", sig))

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//// catching ctx.Done(). timeout of 5 seconds.
	//select {
	//case <-ctx.Done():
	//	log.Warn("server shutdown with timeout of 2 seconds.")
	//}
	//log.Info("Server exiting")
}
