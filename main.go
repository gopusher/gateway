package main

import (
	"flag"
	"runtime"
	"context"
	"time"

	c "gopusher/comet/config"
	"gopusher/comet/contracts"
	"gopusher/comet/connection/websocket"
	"gopusher/comet/service"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"log"
	"fmt"
	"encoding/json"
	"gopusher/comet/discovery"
	"gopusher/comet/rpc"
)

func main() {
	filename, isMonitor := getArgs()
	config := c.NewConfig(*filename)
	runtime.GOMAXPROCS(config.Get("max_proc").MustInt(runtime.NumCPU()))
	cometServiceName := config.Get("comet_service_name").MustString("comet")
	etcdAddr := []string{config.Get("etcd_addr").String()}

	rpcClient := rpc.NewClient(config)

	if *isMonitor {
		discovery.Run(etcdAddr, cometServiceName, rpcClient)
		return
	}
	//todo 信号接管方便平滑重启（目前不处理以后增加

	server := getCometServer(config, rpcClient)

	go server.Run()

	go service.InitRpcServer(server)

	//join cluster
	joinCluster(config, etcdAddr, cometServiceName, server.GetRpcAddr(), server.GetCometAddr())
}

func getCometServer(config *c.Config, rpcClient *rpc.Client) contracts.Server {
	socketProtocol := config.Get("socket_protocol").MustString("websocket")
	switch socketProtocol {
	case "websocket":
		return websocket.NewWebSocketServer(config, rpcClient)
	case "tcp": //暂时不处理
		panic("不支持的通信协议:" + socketProtocol)
	default:
		panic("不支持的通信协议:" + socketProtocol)
	}
}

func getArgs() (filename *string, isMonitor *bool) {
	filename = flag.String("c", "./comet.ini", "set config file path")
	//是否为 monitor 节点
	isMonitor = flag.Bool("m", false, "if running with monitor model")
	flag.Parse()

	return
}

func joinCluster(config *c.Config, etcdServer []string, serviceName string, rpcAddr string, cometAddr string) {
	key := fmt.Sprintf("/%s/%s", serviceName, rpcAddr)
	ttl := 1

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdServer,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	type etcdValue struct {
		Protocol 	string `protocol`
		RpcAddr 	string	`rpc_addr`
		CometAddr	string	`comet_addr`
	}

	body, err := json.Marshal(&etcdValue{
		Protocol: config.Get("socket_protocol").MustString("websocket"),
		RpcAddr: rpcAddr,
		CometAddr: cometAddr,
	})

	log.Println(fmt.Sprintf("rpcAddr: %s, etcdValue: %s, 加入集群成功", rpcAddr, string(body)))

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	var curLeaseId clientv3.LeaseID = 0

	for {
		if curLeaseId == 0 {
			leaseResp, err := lease.Grant(context.TODO(), int64(ttl + 2))
			if err != nil {
				goto SLEEP
			}

			if _, err := kv.Put(context.TODO(), key, string(body), clientv3.WithLease(leaseResp.ID)); err != nil {
				goto SLEEP
			}
			curLeaseId = leaseResp.ID
		} else {
			//log.Printf("keepalive curLeaseId=%d\n", curLeaseId)
			if _, err := lease.KeepAliveOnce(context.TODO(), curLeaseId); err == rpctypes.ErrLeaseNotFound {
				curLeaseId = 0
				continue
			}
		}
	SLEEP:
		time.Sleep(time.Duration(ttl) * time.Second)
	}
}
