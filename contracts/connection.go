package contracts

type Server interface {
	Run()
	SendToConnections(to []string, msg string) ([]string, error)
	GetRpcAddr() string
	GetCometAddr() string
}
