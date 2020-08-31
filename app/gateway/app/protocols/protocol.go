package protocols

import (
	"encoding/json"
)

//Server is comet server interface
type Server interface {
	Protocol() string
	Config() interface{}
	Run() error
	JoinCluster() error
	LeaveCluster() error
	SendToConnections(connections []string, msg string) ([]string, error)
	Broadcast(msg string)
	KickConnections(connections []string)
	KickAllConnections()
	CheckConnectionsOnline(connections []string) []string
	GetAllConnections() []string
	AnyCall(method string, args json.RawMessage) (interface{}, error)
}
