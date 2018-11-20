package configuration

import (
	"os"
	"strings"
)

type MonitorConfig struct {
	EtcdServers 			[]string
	ServiceName				string

	NotificationUrl			string
	NotificationUserAgent	string
}

func GetMonitorConfig() *MonitorConfig {
	etcdServersEnv := os.Getenv("ETCD_SERVER_ADDRESSES")
	if etcdServersEnv == "" {
		panic("need env: ETCD_SERVER_ADDRESSES")
	}
	etcdServers := strings.Split(etcdServersEnv, ",")

	serviceName := os.Getenv("GOPUSHER_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "Gopusher"
	}

	notificationUserAgent := os.Getenv("NOTIFICATION_USER_AGENT")
	if notificationUserAgent == "" {
		notificationUserAgent = "Gopusher 1.0"
	}

	return &MonitorConfig {
		EtcdServers: etcdServers,
		ServiceName: serviceName,
		NotificationUrl: os.Getenv("NOTIFICATION_URL"),
		NotificationUserAgent: notificationUserAgent,
	}
}
