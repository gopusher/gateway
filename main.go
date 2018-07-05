package main

import (
	"flag"
	"runtime"
	c "gopusher/comet/config"
	"gopusher/comet/contracts"
	"gopusher/comet/connection/websocket"
	"gopusher/comet/service"
	"log"
	"fmt"
	"gopusher/comet/discovery"
	"gopusher/comet/rpc"
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
	joinCluster(config, discoveryService, server.GetRpcAddr())
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

func joinCluster(config *c.Config, discoveryService *discovery.Discovery, rpcAddr string) {
	log.Println(fmt.Sprintf("rpcAddr: %s, etcdValue: %s, 加入集群成功", rpcAddr))

	discoveryService.KeepAlive(rpcAddr, rpcAddr)
}

func addComet(rpcClient *rpc.Client) func(string, string) {
	return func(node string, nodeInfo string) {
		fmt.Println("增加节点: " + node)

		if _, err := rpcClient.SuccessRpc("Im", "addCometServer", node, nodeInfo); err != nil {
			color.Red("增加节点失败>> " + err.Error())
		}
	}
}

func delComet(rpcClient *rpc.Client) func(string) {
	return func(node string) {
		fmt.Println("移除节点: " + node)

		if _, err := rpcClient.SuccessRpc("Im", "removeCometServer", node); err != nil {
			color.Red("移除节点>> " + err.Error())
		}
	}
}
