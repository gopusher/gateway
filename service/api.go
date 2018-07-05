package service

import (
	"net/rpc"
	"gopusher/comet/contracts"
	"net"
	"net/rpc/jsonrpc"
	"reflect"
	"encoding/json"
	"errors"
	"github.com/fatih/color"
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

func (s *Server) SendToConnections(body string, reply *string) error {
	//const messageMaxLen = 200
	//if strings.Count(body, "") - 1 > messageMaxLen {
	//	return errors.New(fmt.Sprintf("消息体过长，最大允许长度: %d", messageMaxLen))
	//}

	type Message struct {
		To   	[]string	`json:"to"`	//消息接受者
		Msg 	string		`json:"msg"` //为一个json，里边包含 type 消息类型
		Token	string		`json:"token"` //作为消息发送鉴权
	}

	type Response struct {
		ErrIds	[]string	`json:"error_ids"`
		ErrInfo	string		`json:"error_info"`
	}

	var message Message
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		color.Red("消息体异常, 不能解析 %v %v", body, reflect.TypeOf(body))

		response, _ := json.Marshal(&Response{
			ErrIds: []string{},
			ErrInfo: "msg json marshal error.",
		})
		return errors.New(string(response))
	}

	if message.Token != s.token {
		response, _ := json.Marshal(&Response{
			ErrIds: []string{},
			ErrInfo: "token error.",
		})
		return errors.New(string(response))
	}

	if errIds, err := s.server.SendToConnections(message.To, message.Msg); err != nil {
		response, _ := json.Marshal(&Response{
			ErrIds: errIds,
			ErrInfo: "send failed, " + err.Error(),
		})
		return errors.New(string(response))
	}

	*reply = "消息发送成功"
	return nil
}
