package discovery

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"time"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"fmt"
	"log"
	"gopusher/comet/rpc"
	"github.com/fatih/color"
)

type Discovery struct {
	etcdServer []string
	serviceName string
	rpcClient *rpc.Client
}

func Run(etcdServer []string, serviceName string, rpcClient *rpc.Client) {
	serviceName = fmt.Sprintf("/%s/", serviceName)

	discovery := &Discovery{
		etcdServer: etcdServer,
		serviceName: serviceName,
		rpcClient: rpcClient,
	}

	log.Println("monitor running...")
	discovery.Watch(context.TODO())
}

func (discovery *Discovery) Watch(ctx context.Context) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   discovery.etcdServer,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	var curRevision int64 = 0

	// 先读当前所有孩子, 直到成功为止
	kv := clientv3.NewKV(client)
	for {
		rangeResp, err := kv.Get(context.TODO(), discovery.serviceName, clientv3.WithPrefix())
		if err != nil {
			continue
		}

		for _, kv := range rangeResp.Kvs {
			discovery.addComet(string(kv.Key), string(kv.Value))
		}

		// 从当前版本开始订阅
		curRevision = rangeResp.Header.Revision + 1
		break
	}

	// 监听后续的PUT与DELETE事件
	watcher := clientv3.NewWatcher(client)
	defer watcher.Close()

	watchChan := watcher.Watch(ctx, discovery.serviceName, clientv3.WithPrefix(), clientv3.WithRev(curRevision))
	for watchResp := range watchChan { // if ctx is Done, for loop will break
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				//fmt.Println("PUT事件: " + string(event.Kv.Key) + ">>" + string(event.Kv.Value))
				discovery.addComet(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				//fmt.Println("DELETE事件: " + string(event.Kv.Key) + ">>" + string(event.Kv.Value))
				discovery.delComet(string(event.Kv.Key))
			}
		}
	}
}

func (discovery *Discovery) addComet(nodeName string, value string) {
	fmt.Println("增加节点: " + nodeName)

	if _, err := discovery.rpcClient.SuccessRpc("Im", "addCometServer", nodeName, value); err != nil {
		color.Red(err.Error())
	}
}

func (discovery *Discovery) delComet(nodeName string) {
	fmt.Println("移除节点: " + nodeName)

	if _, err := discovery.rpcClient.SuccessRpc("Im", "removeCometServer", nodeName); err != nil {
		color.Red(err.Error())
	}
}
