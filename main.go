package main

import (
	"flag"
	"runtime"
	c "github.com/gopusher/comet/config"
	"github.com/gopusher/comet/contracts"
	"github.com/gopusher/comet/connection/websocket"
	"github.com/gopusher/comet/service"
	"fmt"
	"github.com/gopusher/comet/discovery"
	"github.com/gopusher/comet/rpc"
	"github.com/fatih/color"
)

func getArgs() (filename *string, isMonitor *bool) {
	filename = flag.String("c", "./comet.ini", "set config file path")
	//是否为 monitor 节点
	isMonitor = flag.Bool("m", false, "if running with monitor model")
	flag.Parse()

	return
}

func main() {
	filename, isMonitor := getArgs()

	config := c.NewConfig(*filename)
	runtime.GOMAXPROCS(config.Get("max_proc").MustInt(runtime.NumCPU()))
	cometServiceName := config.Get("comet_service_name").MustString("comet")
	etcdAddr := []string{config.Get("etcd_addr").String()}

	rpcClient := rpc.NewClient(config.Get("rpc_api_url").String(), config.Get("rpc_user_agent").String())

	discoveryService := discovery.NewDiscovery(etcdAddr, cometServiceName)
	if *isMonitor {
		discoveryService.Watch(addComet(rpcClient), delComet(rpcClient))
		return
	}
	//todo 信号接管方便平滑重启（目前不处理以后增加

	server := getCometServer(config, rpcClient)

	go server.Run()

	go service.InitRpcServer(server, config.Get("comet_rpc_token").MustString("token"))

	//join cluster
	joinCluster(discoveryService, server.GetRpcAddr())
}

func getCometServer(config *c.Config, rpcClient *rpc.Client) contracts.Server {
	socketProtocol := config.Get("socket_protocol").MustString("ws")
	switch socketProtocol {
	case "ws":
		fallthrough
	case "wss":
		return websocket.NewWebSocketServer(config, rpcClient)
	case "tcp": //暂时不处理
		panic("不支持的通信协议:" + socketProtocol)
	default:
		panic("不支持的通信协议:" + socketProtocol)
	}
}

func joinCluster(discoveryService *discovery.Discovery, rpcAddr string) {
	discoveryService.KeepAlive(rpcAddr)
}

func addComet(rpcClient *rpc.Client) func(string, string) {
	return func(node string, revision string) {
		fmt.Printf("增加节点: node: %s, revision: %s \n", node, revision)

		if _, err := rpcClient.SuccessRpc("Im", "addCometServer", node, revision); err != nil {
			color.Red("增加节点失败: " + err.Error())
		}
	}
}

func delComet(rpcClient *rpc.Client) func(string, string) {
	return func(node string, revision string) {
		fmt.Printf("移除节点: node: %s, revision: %s \n", node, revision)

		if _, err := rpcClient.SuccessRpc("Im", "removeCometServer", node, revision); err != nil {
			color.Red("移除节点失败: " + err.Error())
		}
	}
}
