package comet

import (
	"github.com/gopusher/gateway/notification"
	"github.com/gopusher/gateway/configuration"
	"github.com/gopusher/gateway/contracts"
	"github.com/gopusher/gateway/connection/websocket"
	"github.com/gopusher/gateway/api"
	"github.com/gopusher/gateway/discovery"
)

func Run() {
	config := configuration.GetCometConfig()

	server := getCometServer(config)

	go server.Run()

	go api.InitRpcServer(server, config)

	joinCluster(config)
}

func getCometServer(config *configuration.CometConfig) contracts.Server {
	rpc := notification.NewRpc(config.NotificationUrl, config.NotificationUserAgent)

	switch config.SocketProtocol {
	case "ws":
		fallthrough
	case "wss":
		return websocket.NewWebSocketServer(config, rpc)
	case "tcp": //暂时不处理
		panic("Unsupported protocol: " + config.SocketProtocol)
	default:
		panic("Unsupported protocol: " + config.SocketProtocol)
	}
}

func joinCluster(config *configuration.CometConfig) {
	node := config.SocketAddress + "-" + config.GatewayApiAddress + "-" + config.GatewayApiToken
	discovery.NewDiscovery(config.EtcdServers, config.ServiceName).KeepAlive(node)
}
