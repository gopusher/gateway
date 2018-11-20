package api

import (
	"net/rpc"
	"github.com/gopusher/gateway/contracts"
	"net"
	"net/rpc/jsonrpc"
	"github.com/gopusher/gateway/log"
	"github.com/gopusher/gateway/configuration"
	"encoding/json"
	)

type Server struct {
	server contracts.Server
	token string
}

func InitRpcServer(server contracts.Server, config *configuration.CometConfig) {
	rpc.Register(&Server{
		server: server,
		token: config.GatewayApiToken,
	})

	listener, err := net.Listen("tcp", config.GatewayApiPort)
	if err != nil {
		panic("Gateway api server run failed, error: %s" + err.Error())
	}

	log.Info("Gateway api server start running: %s", config.GatewayApiAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}

type TokenMessage struct {
	Token			string		`json:"token"` 		//作为消息发送鉴权
}

type ConnectionsMessage struct {
	Connections		[]string	`json:"connections"`	//消息接受者
	TokenMessage
}

type Response struct {
	Connections		[]string	`json:"connections"`	//消息接受者
	Error			string		`json:"error"`
}

func (s *Server) checkToken(token string) string {
	if token != s.token {
		response, _ := json.Marshal(&Response{
			Connections: 	[]string{},
			Error:			"error token",
		})

		return string(response)
	}

	return ""
}

func (s *Server) success(connections []string) string {
	if connections == nil {
		connections = []string{}
	}

	response, _ := json.Marshal(&Response{
		Connections:	connections,
		Error:			"",
	})

	return string(response)
}

func (s *Server) failure(connections []string, err string) string {
	if connections == nil {
		connections = []string{}
	}

	response, _ := json.Marshal(&Response{
		Connections:	connections,
		Error:			err,
	})

	return string(response)
}
