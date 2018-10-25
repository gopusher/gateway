package contracts

type Server interface {
	Run()
	Broadcast(msg string)
	SendToConnections(connections []string, msg string) ([]string, error)
	GetRpcAddr() string
	GetRpcPort() string
	KickConnections(connections []string)
}
