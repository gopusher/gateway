package service

import (
	"net/rpc"
	"gopusher/comet/contracts"
	"net"
	"net/rpc/jsonrpc"
	"encoding/json"
	"errors"
)

type Server struct {
	server contracts.Server
	token string
}

func InitRpcServer(server contracts.Server, token string) {
	rpc.Register(&Server{
		server: server,
		token: token,
	})
	listener, err := net.Listen("tcp", server.GetRpcAddr())
	if err != nil {
		panic("rpc服务初始化失败, " + err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		//新协程来处理--json
		go jsonrpc.ServeConn(conn)
	}
}

type Message struct {
	Connections		[]string	`json"connections"`	//消息接受者
	Msg 			string		`json"msg"` 		//为一个json，里边包含 type 消息类型
	Token			string		`json"token"` 		//作为消息发送鉴权
}

func (s *Server) SendToConnections(message *Message, reply *string) error {
//token string, connections []string, msg string
	type Response struct {
		ErrIds	[]string	`json:"ids"`
		ErrInfo	string		`json:"msg"`
	}

	if message.Token != s.token {
		response, _ := json.Marshal(&Response{
			ErrIds: []string{},
			ErrInfo: "token error.",
		})
		return errors.New(string(response))
	}

	if len(message.Connections) == 0 {
		response, _ := json.Marshal(&Response{
			ErrIds: []string{},
			ErrInfo: "empty connections.",
		})
		return errors.New(string(response))
	}

	if errIds, err := s.server.SendToConnections(message.Connections, message.Msg); err != nil {
		response, _ := json.Marshal(&Response{
			ErrIds: errIds,
			ErrInfo: "send failed, " + err.Error(),
		})
		return errors.New(string(response))
	}

	*reply = "ok"
	return nil
}

func (s *Server) KickConnections(connections []string, reply *string) error {
	go s.server.KickConnections(connections)

	*reply = "ok"
	return nil
}
