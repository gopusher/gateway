package main

import (
	"flag"
	"runtime"
	"context"
	"time"

	c "gopusher/comet/config"
	"gopusher/comet/contracts"
	"gopusher/comet/connection/websocket"
	"gopusher/comet/rpc"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"log"
	"fmt"
	"encoding/json"
)

func main() {
	filename := getArgs()
	config := c.NewConfig(*filename)
	runtime.GOMAXPROCS(config.Get("max_proc").MustInt(runtime.NumCPU()))
	//todo 信号接管方便平滑重启（目前不处理以后增加）

	server := getCometServer(config)

	go server.Run()

	go rpc.InitRpcServer(server)

	//join cluster
	nodeName := config.Get("node_name").MustString("")
	if nodeName == "" {
		panic("node_name 为空")
	}
	joinCluster([]string{config.Get("etcd_addr").String()}, "comet", nodeName, server.GetRpcAddr(), server.GetCometAddr())
}

func getCometServer(config *c.Config) contracts.Server {
	socketProtocol := config.Get("socket_protocol").MustString("websocket")
	switch socketProtocol {
	case "websocket":
		return websocket.NewWebSocketServer(config)
	case "tcp": //暂时不处理
		panic("不支持的通信协议:" + socketProtocol)
	default:
		panic("不支持的通信协议:" + socketProtocol)
	}
}

func getArgs() (filename *string) {
	filename = flag.String("c", "./comet.ini", "set config file path")
	flag.Parse()

	return filename
}

func joinCluster(etcdServer []string, serviceName string, nodeName string, rpcAddr string, cometAddr string) {
	key := fmt.Sprintf("/%s/%s", serviceName, nodeName)
	ttl := 1

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdServer,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	type etcdValue struct {
		RpcAddr 	string	`rpc_addr`
		CometAddr	string	`comet_addr`
	}

	body, err := json.Marshal(&etcdValue{
		RpcAddr: rpcAddr,
		CometAddr: cometAddr,
	})

	log.Println(nodeName + " 加入集群成功" + string(body))

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
