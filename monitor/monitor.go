package monitor

import (
	"github.com/gopusher/gateway/discovery"
	"github.com/gopusher/gateway/notification"
	"github.com/gopusher/gateway/log"
	"github.com/gopusher/gateway/configuration"
)

var rpc *notification.Client

func Run() {
	config := configuration.GetMonitorConfig()

	rpc = notification.NewRpc(config.NotificationUrl, config.NotificationUserAgent)

	discoveryService := discovery.NewDiscovery(config.EtcdServers, config.ServiceName)
	discoveryService.Watch(addComet(), removeComet())
}

func addComet() func(string, string) {
	return func(node string, revision string) {
		if _, err := rpc.Call("AddServer", node, revision); err != nil {
			log.Error("notification failed: add node, node: %s, revision: %s, error: %s", node, revision, err.Error())
			return
		}

		log.Info("add node: node: %s, revision: %s success", node, revision)
	}
}

func removeComet() func(string, string) {
	return func(node string, revision string) {
		if _, err := rpc.Call("RemoveServer", node, revision); err != nil {
			log.Error("notification failed: remove node, node: %s, revision: %s, error: %s", node, revision, err.Error())
			return
		}

		log.Info("remove node: node: %s, revision: %s success", node, revision)
	}
}
