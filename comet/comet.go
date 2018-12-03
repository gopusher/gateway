package comet

import (
	"github.com/gopusher/gateway/configuration"
	"github.com/gopusher/gateway/contracts"
	"github.com/gopusher/gateway/connection/websocket"
	"github.com/gopusher/gateway/api"
)

func Run() {
	config := configuration.GetCometConfig()

	server := getCometServer(config)

	go server.Run()

	go server.JoinCluster()

	api.InitRpcServer(server, config)
}

func getCometServer(config *configuration.CometConfig) contracts.Server {
	switch config.SocketProtocol {
	case "ws":
		fallthrough
	case "wss":
		return websocket.NewWebSocketServer(config)
	case "tcp": //暂时不处理
		panic("Unsupported protocol: " + config.SocketProtocol)
	default:
		panic("Unsupported protocol: " + config.SocketProtocol)
	}
}
