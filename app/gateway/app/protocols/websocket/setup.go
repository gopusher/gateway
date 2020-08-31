package websocket

import (
	"github.com/gopusher/gateway/app/gateway/app/protocols"
)

const (
	protocol = "websocket"
)

func init() {
	protocols.RegisterType(protocol, setup)
}

func setup() (protocols.Server, error) {
	server := newServer()

	return server, nil
}
