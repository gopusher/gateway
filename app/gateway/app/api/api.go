package api

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gopusher/gateway/app/gateway/app/protocols"
	"github.com/gopusher/gateway/pkg/log"
	"go.uber.org/zap"
)

type Config struct {
	Address string `mapstructure:"address" validate:"required"`
	Token   string `mapstructure:"token"`
}

//Server is api server
type Server struct {
	node   string
	server protocols.Server
	config *Config
}

//InitRpcServer init rpc server
func InitRpcServer(node string, server protocols.Server, config *Config) {
	if err := rpc.Register(&Server{
		node:   node,
		server: server,
		config: config,
	}); err != nil {
		log.Panic("Gateway api server run failed, error: "+err.Error(), zap.Error(err))
	}

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Panic("Gateway api server run failed, error: "+err.Error(), zap.Error(err))
	}

	log.Info(fmt.Sprintf(
		"Gateway api server start running, Node: %s, gateway api address: %s, token: %s",
		node, config.Address, config.Token,
	))
	defer func() {
		if err := listener.Close(); err != nil {
			log.Panic("Gateway api server run failed, error: "+err.Error(), zap.Error(err))
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}

//TokenMessage is token type message
type TokenMessage struct {
	Token string `json:"token"` //作为消息发送鉴权
}

//ConnectionsMessage is a connection type message
type ConnectionsMessage struct {
	Connections []string `json:"connections"` //消息接受者
	TokenMessage
}

//Response is api response
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

const (
	SuccessCode = 20000

	ParametersError = 40000

	TokenErrorCode = 40100

	ServerErrorCode = 50000

	SendToConnectionsErrorCode = 50001
)

func (s *Server) success(data interface{}) Response {
	return Response{
		Code:  SuccessCode,
		Data:  data,
		Error: "",
	}
}

func (s *Server) failure(code int, err string) Response {
	return Response{
		Code:  code,
		Data:  "",
		Error: err,
	}
}

func (s *Server) failureWithData(code int, data interface{}, err string) Response {
	return Response{
		Code:  code,
		Data:  data,
		Error: err,
	}
}

func (s *Server) checkToken(token string) *Response {
	if token != s.config.Token {
		resp := s.failure(TokenErrorCode, "error rpc token")
		return &resp
	}

	return nil
}
