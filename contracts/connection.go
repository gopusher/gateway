package contracts

type Server interface {
	Run()
	SendToConnections(connections []string, msg string) ([]string, error)
	GetRpcAddr() string
	KickConnections(connections []string) error
}
