package discovery

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type Discovery struct {
	etcdServer []string
	serviceName string
}

func NewDiscovery(etcdServer []string, serviceName string) *Discovery {
	return &Discovery{
		etcdServer: etcdServer,
		serviceName: "/" + serviceName,
	}
}

func (discovery *Discovery) getClient() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   discovery.etcdServer,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	return client
}

func (discovery *Discovery) KeepAlive(node string) {
	key := fmt.Sprintf("%s/%s", discovery.serviceName, node)
	ttl := 1
	client := discovery.getClient()

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	var curLeaseId clientv3.LeaseID = 0

	leaseResp, err := lease.Grant(context.TODO(), int64(ttl + 1))
	if err != nil {
		panic(err)
	}

	if _, err := kv.Put(context.TODO(), key, fmt.Sprintf("%d", leaseResp.GetRevision() + 1), clientv3.WithLease(leaseResp.ID)); err != nil {
		panic(err)
	}

	curLeaseId = leaseResp.ID

	for {
		time.Sleep(time.Duration(ttl) * time.Second)

		if _, err := lease.KeepAliveOnce(context.TODO(), curLeaseId); err == rpctypes.ErrLeaseNotFound {
			panic(err)
		}
	}
}

func (discovery *Discovery) Watch(addClient func(string, string), removeClient func(string, string)) {
	client := discovery.getClient()

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
			addClient(string(kv.Key), string(kv.Value))
		}

		// 从当前版本开始订阅
		curRevision = rangeResp.Header.Revision + 1
		break
	}

	// 监听后续的PUT与DELETE事件
	watcher := clientv3.NewWatcher(client)
	defer watcher.Close()

	watchChan := watcher.Watch(context.TODO(), discovery.serviceName, clientv3.WithPrefix(), clientv3.WithRev(curRevision))
	for watchResp := range watchChan { // if ctx is Done, for loop will break
		for _, event := range watchResp.Events {
			//fmt.Printf("watch 事件: %d, %d \n", watchResp.Header.GetRevision(), event.Kv.Version)
			switch event.Type {
			case mvccpb.PUT:
				//fmt.Println("PUT事件: " + string(event.Kv.Key) + ">>" + string(event.Kv.Value))
				addClient(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				//fmt.Println("DELETE事件: " + string(event.Kv.Key) + ">>" + string(event.Kv.Value))
				removeClient(string(event.Kv.Key), fmt.Sprintf("%d", watchResp.Header.GetRevision()))
			}
		}
	}
}
