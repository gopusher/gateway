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
	"bytes"
	"io/ioutil"
)

type Server struct {
	config *config.Config
	apiAddr string
	wsAddr string
	rpcAddr string
	upgrader websocket.Upgrader
	register chan *Client
	unregister chan *Client
	clients map[string]*Client
}

func NewWebSocketServer(config *config.Config) *Server {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	apiAddr := config.Get("api_addr").String()
	wsAddr := config.Get("websocket_port").MustString(":8900")
	rpcAddr := config.Get("rpc_addr").MustString("127.0.0.1:8901")

	return &Server{
		config: config,
		apiAddr: apiAddr,
		wsAddr: wsAddr,
		rpcAddr: rpcAddr,
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

func (s *Server) GetCometAddr() string {
	websocketProtocol := s.config.Get("websocket_protocol").MustString("ws")
	websocketHost := s.config.Get("websocket_host").MustString("127.0.0.1")
	return fmt.Sprintf("%s://%s%s", websocketProtocol, websocketHost, s.wsAddr)
}

// 启动 websocket server
func (s *Server) initWsServer() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ws", s.serveWs)

	log.Println("[info] websocket server start running: " + s.wsAddr)
	websocketProtocol := s.config.Get("websocket_protocol").MustString("ws")
	if websocketProtocol == "ws" {
		if err := http.ListenAndServe(s.wsAddr, serverMux); err != nil {
			log.Fatal("服务启动失败:" + err.Error())
			panic(err)
		}
	} else {
		wssCertPem := s.config.Get("wss_cert_pem").String()
		wssKeyPem := s.config.Get("wss_key_pem").MustString("ws")
		if err := http.ListenAndServeTLS(s.wsAddr, wssCertPem, wssKeyPem, serverMux); err != nil {
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
			if _, err := s.successRpc("Im", "online", client.connId, client.info, s.rpcAddr); err != nil {
				color.Red(err.Error())
			}
		case client := <-s.unregister:
			log.Println("[info] 断开连接，connId:" + client.connId)
			//上报给 router api 服务
			if _, err := s.successRpc("Im", "offline", client.connId, client.info, s.rpcAddr); err != nil {
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

func (s Server) doRpc(class string, method string, args ...interface{}) (string, error) {
	type RpcBody struct {
		ClassName	string	`json:"class"`
		MethodName 	string	`json:"method"`
		Args		[]interface{}	`json:"args"`
	}
	body, err := json.Marshal(&RpcBody{
		ClassName: class,
		MethodName: method,
		Args: args,
	})

	apiUrl := s.config.Get("rpc_api_url").String()
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("请求失败:" + err.Error())
	}
	req.Header.Set("User-Agent", s.config.Get("rpc_user_agent").String())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		ret, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(ret), nil
	}

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("请求异常%d: %v", resp.StatusCode, err.Error()))
	}

	return "", errors.New(fmt.Sprintf("请求异常%d: %v", resp.StatusCode, string(ret)))
}

func (s Server) successRpc(class string, method string, args ...interface{}) (string, error) {
	ret, err := s.doRpc(class, method, args...)
	if err != nil {
		return "", err
	}

	type RetInfo struct {
		Code int `code`
		Data interface{} `data`
		Error string `error`
	}

	var retInfo RetInfo
	if err := json.Unmarshal([]byte(ret), &retInfo); err != nil {
		color.Red("消息体异常, 不能解析 %v %v", ret, reflect.TypeOf(ret))
		return "", errors.New("消息体异常, 不能解析")
	}

	if retInfo.Code != 0 {
		return "", errors.New(retInfo.Error)
	}

	return ret, nil
}

func (s Server) serveWs(w http.ResponseWriter, r *http.Request) {
	//检查是否是有效连接
	tokenInfo, err := s.checkToken(r.URL.Query())
	if err != nil {
		s.responseWsUnauthorized(w)
		return
	}

	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn: c,
		send: make(chan []byte, 256),
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

	if _, err := s.successRpc("Im", "checkToken", tokenInfo.ConnId, tokenInfo.Token, tokenInfo.Info, s.GetCometAddr()); err != nil {
		color.Red(err.Error())
		return nil, errors.New("授权失败" + err.Error())
	}

	return &tokenInfo, nil
}

func (s *Server) SendToConnections(to []string, msg string) error {
	for _, id := range to {
		if err := s.SendToConnection(id, msg); err != nil {
			//todo 优化：将失败的放入数组返回或则记录日志和原因等
		}
	}

	return nil
}

func (s *Server) SendToConnection(to string, msg string) error {
	if client, ok := s.clients[to]; ok {
		select {
		case client.send <- []byte(msg):
			log.Println("[info] SendToConnection " + to + ": " + msg)
			return nil
		default:
			close(client.send)
			delete(s.clients, to)
			color.Red("发送消息失败, to: %s", to)
			return errors.New(fmt.Sprintf("发送消息失败, to %s", to))
		}
	}

	color.Red("发送消息失败, 客户端不在维护中, to: %s", to)
	return errors.New(fmt.Sprintf("发送消息失败, 客户端不在维护中, to %s", to))
}
