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
	gatewayApiAddress := os.Getenv("GATEWAY_API_ADDRESS")
	gatewayApiAddressSlice := strings.Split(gatewayApiAddress, ":")
	if len(gatewayApiAddressSlice) != 2 {
		panic("error env: GATEWAY_API_ADDRESS")
	}

	nodeId := gatewayApiAddress + ":" + strconv.FormatInt(time.Now().UnixNano(), 10)

	notificationUserAgent := os.Getenv("NOTIFICATION_USER_AGENT")
	if notificationUserAgent == "" {
		notificationUserAgent = "Gopusher 1.0"
	}

	return &CometConfig {
		NodeId: nodeId,

		SocketProtocol: os.Getenv("SOCKET_PROTOCOL"),
		SocketPort: ":" + os.Getenv("SOCKET_PORT"),
		SocketCertFile: os.Getenv("SOCKET_CERT_FILE"),
		SocketKeyFile: os.Getenv("SOCKET_KEY_FILE"),

		GatewayApiAddress: gatewayApiAddress,
		GatewayApiPort: ":" + gatewayApiAddressSlice[1:][0],
		GatewayApiToken: os.Getenv("GATEWAY_API_TOKEN"),

		NotificationUrl: os.Getenv("NOTIFICATION_URL"),
		NotificationUserAgent: notificationUserAgent,
	}
}
