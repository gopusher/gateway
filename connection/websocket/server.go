package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"gopusher/comet/config"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"encoding/json"
	"reflect"
	"gopusher/comet/rpc"
)

type Server struct {
	config *config.Config
	wsPort string
	rpcAddr string
	rpcClient *rpc.Client
	upgrader websocket.Upgrader
	register chan *Client
	unregister chan *Client
	clients map[string]*Client
}

func NewWebSocketServer(config *config.Config, rpcClient *rpc.Client) *Server {
	var upgrader = websocket.Upgrader{ //todo 搞成配置
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsPort := config.Get("websocket_port").MustString(":8900")
	rpcAddr := config.Get("rpc_addr").MustString("127.0.0.1:8901")

	return &Server{
		config: config,
		wsPort: wsPort,
		rpcAddr: rpcAddr,
		rpcClient: rpcClient,
		upgrader: upgrader,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[string]*Client),
	}
}

func (s *Server) Run() {
	go s.handleClients()

	s.initWsServer()
}

func (s *Server) GetRpcAddr() string {
	return s.rpcAddr
}

// 启动 websocket server
func (s *Server) initWsServer() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ws", s.serveWs)

	log.Println("[info] websocket server start running: " + s.wsPort)
	websocketProtocol := s.config.Get("socket_protocol").MustString("ws")
	if websocketProtocol == "wss" {
		wssCertPem := s.config.Get("wss_cert_pem").String()
		wssKeyPem := s.config.Get("wss_key_pem").MustString("ws")
		if err := http.ListenAndServeTLS(s.wsPort, wssCertPem, wssKeyPem, serverMux); err != nil {
			log.Fatal("服务启动失败:" + err.Error())
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(s.wsPort, serverMux); err != nil {
			log.Fatal("服务启动失败:" + err.Error())
			panic(err)
		}
	}
}

func (s *Server) handleClients() {
	log.Println("[info] handle clients")
	for {
		select {
		case client := <-s.register:
			log.Println("[info] 注册客户端, connId: " + client.connId)

			s.clients[client.connId] = client

			//上报给 router api 服务
			if _, err := s.rpcClient.SuccessRpc("Im", "online", client.connId, client.info, s.rpcAddr); err != nil {
				color.Red(err.Error())
			}
		case client := <-s.unregister:
			log.Println("[info] 断开连接，connId:" + client.connId)
			//上报给 router api 服务
			if _, err := s.rpcClient.SuccessRpc("Im", "offline", client.connId, client.info, s.rpcAddr); err != nil {
				color.Red(err.Error())
			}

			//关闭客户端连接
			if _, ok := s.clients[client.connId]; ok {
				delete(s.clients, client.connId)
				client.Close()
			}
		}
	}
}

func (s Server) serveWs(w http.ResponseWriter, r *http.Request) {
	//检查是否是有效连接
	tokenInfo, err := s.checkToken(r.URL.Query())
	if err != nil {
		s.responseWsUnauthorized(w)
		return
	}

	//存在相同connId客户端连接
	if _, ok := s.clients[tokenInfo.ConnId]; ok {
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn: c,
		send: make(chan []byte, 1024), //todo 搞成配置
		connId: tokenInfo.ConnId,
		info: tokenInfo.Info,
		server: s,
	}

	s.register <- client

	go client.Write()
	go client.Read()
}

func (s Server) responseWsUnauthorized(w http.ResponseWriter) { //todo 移动到 message 中
	w.Header().Set("Sec-Websocket-Version", "13")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

type TokenInfo struct {
	ConnId	string	`json:"conn_id"` //唯一分配的conn_id
	Token	string	`json:"token"` //授权token
	Info	string	`json:"info"`	//一些其他信息 json 串
}

func (s Server) checkToken(query map[string][]string) (*TokenInfo, error) {
	if t, ok := query["t"]; !ok || len(t) < 1 {
		return nil, errors.New("确实参数")
	}
	t := query["t"][0]

	var tokenInfo TokenInfo
	if err := json.Unmarshal([]byte(t), &tokenInfo); err != nil {
		color.Red("消息体异常, 不能解析 %v %v", t, reflect.TypeOf(t))
		return nil, errors.New("消息体异常, 不能解析")
	}

	if _, err := s.rpcClient.SuccessRpc("Im", "checkToken", tokenInfo.ConnId, tokenInfo.Token, tokenInfo.Info); err != nil {
		color.Red(err.Error())
		return nil, errors.New("授权失败" + err.Error())
	}

	return &tokenInfo, nil
}

func (s *Server) SendToConnections(to []string, msg string) ([]string, error) {
	var errIds []string
	for _, id := range to {
		if err := s.SendToConnection(id, msg); err != nil {
			errIds = append(errIds, id)
		}
	}
	if len(errIds) > 0 {
		return errIds, errors.New("存在发送失败的消息")
	}

	return []string{}, nil
}

func (s *Server) SendToConnection(to string, msg string) error {
	if client, ok := s.clients[to]; ok {
		select {
		case client.send <- []byte(msg):
			// log.Println("[info] SendToConnection " + to + ": " + msg)
			return nil
		default:
			delete(s.clients, to)
			close(client.send) //是否需要 关闭 chan 的时候，发送完毕所有的chan再关闭连接 ??
			//client.Close()
			color.Red("发送消息失败, to: %s", to)
			return errors.New(fmt.Sprintf("发送消息失败, to %s", to))
		}
	}

	color.Red("发送消息失败, 客户端不在维护中, to: %s", to)
	return errors.New(fmt.Sprintf("发送消息失败, 客户端不在维护中, to %s", to))
}
