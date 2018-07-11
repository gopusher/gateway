package contracts

type Server interface {
	Run()
	SendToConnections(connections []string, msg string) ([]string, error)
	GetRpcAddr() string
	GetRpcPort() string
	KickConnections(connections []string) error
}
