package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"github.com/gopusher/comet/config"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gopusher/comet/rpc"
)

type Server struct {
	config *config.Config
	wsPort string
	rpcAddr string
	rpcPort string
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

	wsPort := ":" + config.Get("websocket_port").MustString("8900")
	rpcPort := ":" + config.Get("comet_rpc_port").MustString("8901")
	rpcAddr := config.Get("comet_rpc_addr").MustString("127.0.0.1")

	return &Server{
		config: config,
		wsPort: wsPort,
		rpcPort: rpcPort,
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
	return s.rpcAddr + s.rpcPort
}

func (s *Server) GetRpcPort() string {
	return s.rpcPort
}

// 启动 websocket server
func (s *Server) initWsServer() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ws", s.serveWs)

	log.Println("[info] websocket server start running " + s.wsPort)
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
	for {
		select {
		case client := <-s.register:
			log.Println("[info] 注册客户端, connId: " + client.connId)

			s.clients[client.connId] = client

			//上报给 router api 服务
			if _, err := s.rpcClient.SuccessRpc("Im", "online", client.connId, client.info, s.GetRpcAddr()); err != nil {
				color.Red(err.Error())
			}
		case client := <-s.unregister:
			log.Println("[info] 断开连接，connId:" + client.connId)
			//上报给 router api 服务
			if _, err := s.rpcClient.SuccessRpc("Im", "offline", client.connId, client.info, s.GetRpcAddr()); err != nil {
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
	connId, clientInfo, err := s.checkToken(r.URL.Query())
	if err != nil {
		s.responseWsUnauthorized(w)
		return
	}

	//存在相同connId客户端连接
	if _, ok := s.clients[connId]; ok {
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
		connId: connId,
		info: clientInfo,
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

func (s Server) checkToken(query map[string][]string) (string, string, error) {
	if c, ok := query["c"]; !ok || len(c) < 1 {
		return "", "", errors.New("缺少参数c")
	}
	if t, ok := query["t"]; !ok || len(t) < 1 {
		return "", "", errors.New("缺少参数t")
	}
	if i, ok := query["i"]; !ok || len(i) < 1 {
		return "", "", errors.New("缺少参数i")
	}

	connId := query["c"][0]
	token := query["t"][0]
	clientInfo := query["i"][0]

	if _, err := s.rpcClient.SuccessRpc("Im", "checkToken", connId, token, clientInfo, s.GetRpcAddr()); err != nil {
		color.Red(err.Error())
		return "", "", errors.New("授权失败" + err.Error())
	}

	return connId, clientInfo, nil
}

func (s *Server) SendToConnections(connections []string, msg string) ([]string, error) {
	var errIds []string
	for _, connId := range connections {
		if err := s.SendToConnection(connId, msg); err != nil {
			errIds = append(errIds, connId)
		}
	}
	if len(errIds) > 0 {
		return errIds, errors.New("存在发送失败的消息")
	}

	return []string{}, nil
}

func (s *Server) SendToConnection(connId string, msg string) error {
	if client, ok := s.clients[connId]; ok {
		select {
		case client.send <- []byte(msg):
			// log.Println("[info] SendToConnection " + to + ": " + msg)
			return nil
		default:
			delete(s.clients, connId)
			close(client.send) //是否需要 关闭 chan 的时候，发送完毕所有的chan再关闭连接 ??
			//client.Close()
			color.Red("发送消息失败, to: %s", connId)
			return errors.New(fmt.Sprintf("发送消息失败, to %s", connId))
		}
	}

	color.Red("发送消息失败, 客户端不在维护中, to: %s", connId)
	return errors.New(fmt.Sprintf("发送消息失败, 客户端不在维护中, to %s", connId))
}

func (s *Server) KickConnections(connections []string) error {
	for _, connId := range connections {
		if client, ok := s.clients[connId]; ok {
			s.unregister <- client
		}
	}

	return nil
}
