package configuration

import (
    "os"
    "strings"
	"strconv"
	"time"
)

type CometConfig struct {
	NodeId 					string
	SocketProtocol			string
	SocketAddress			string
	SocketPort				string
	SocketCertFile			string
	SocketKeyFile			string

	NotificationUrl			string
	NotificationUserAgent	string

	GatewayApiAddress		string
	GatewayApiPort			string
	GatewayApiToken			string
}

func GetCometConfig() *CometConfig {
	socketAddress := os.Getenv("SOCKET_ADDRESS")
	socketAddressSlice := strings.Split(socketAddress, ":")
	if len(socketAddressSlice) != 3 {
		panic("error env: SOCKET_ADDRESS")
	}

	gatewayApiAddress := os.Getenv("GATEWAY_API_ADDRESS")
	gatewayApiAddressSlice := strings.Split(gatewayApiAddress, ":")
	if len(gatewayApiAddressSlice) != 2 {
		panic("error env: SOCKET_ADDRESS")
	}

	nodeId := gatewayApiAddress + ":" + strconv.FormatInt(time.Now().UnixNano(), 10)

	notificationUserAgent := os.Getenv("NOTIFICATION_USER_AGENT")
	if notificationUserAgent == "" {
		notificationUserAgent = "Gopusher 1.0"
	}

	return &CometConfig {
		NodeId: nodeId,

		SocketProtocol: os.Getenv("SOCKET_PROTOCOL"),
		SocketAddress: socketAddress,
		SocketPort: ":" + socketAddressSlice[2:][0],
		SocketCertFile: os.Getenv("SOCKET_CERT_FILE"),
		SocketKeyFile: os.Getenv("SOCKET_KEY_FILE"),

		GatewayApiAddress: gatewayApiAddress,
		GatewayApiPort: ":" + gatewayApiAddressSlice[1:][0],
		GatewayApiToken: os.Getenv("GATEWAY_API_TOKEN"),

		NotificationUrl: os.Getenv("NOTIFICATION_URL"),
		NotificationUserAgent: notificationUserAgent,
	}
}
