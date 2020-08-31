package websocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gopusher/gateway/app/gateway/app/cfg"
	"github.com/gopusher/gateway/pkg/helper"
	"github.com/gopusher/gateway/pkg/log"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server struct {
	config    *Config
	signature *helper.Signature
	node      string

	upgrader   websocket.Upgrader
	register   chan *Client
	unregister chan *Client
	clients    sync.Map //map[string]*Client
}

func newServer() *Server {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	s := &Server{
		config:    defaultConfig,
		signature: helper.NewSignature([]byte(defaultConfig.AppSecret)),
		node:      cfg.Config.Node,

		upgrader:   upgrader,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		//clients:    make(map[string]*Client),
	}
	// parse config
	if err := cfg.Config.Unpack(s); err != nil {
		log.Panic("Unpack panic", zap.Error(err))
	}

	log.Info("newServer, server config for " + protocol)
	log.Debug("defaultConfig:" + helper.ToJsonString(s.config))

	return s
}

func (s Server) Protocol() string {
	return protocol
}

func (s Server) Config() interface{} {
	return s.config
}

func (s Server) Run() error {
	go s.handleClients()

	go s.initWsServer()

	return nil
}

func (s *Server) handleClients() {
	for {
		select {
		case client := <-s.register:
			s.clients.Store(client.clientId, client)
			//s.clients[client.clientId] = client

			//notify router api server
			//todo
			//if _, err := s.rpc.Call("Online", client.clientId, s.config.Node); err != nil {
			//	log.Error("Client online notification failed: %s" + err.Error(), zap.Error(err))
			//}
		case client := <-s.unregister:
			if _, ok := s.clients.LoadAndDelete(client.clientId); ok {
				//delete(s.clients, client.clientId)
				close(client.send)
				client.Close()

				//notify router api server
				//todo
				//if _, err := s.rpc.Call("Offline", client.connId, s.config.Node); err != nil {
				//	log.Error("Client online notification failed: %s", err.Error())
				//}
			}
		}
	}
}

// run websocket server
func (s *Server) initWsServer() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", s.serveWs)

	log.Info("Websocket server start running, port: " + s.config.Address)
	if s.config.Ssl {
		if err := http.ListenAndServeTLS(s.config.Address, s.config.SslCertFile, s.config.SslKeyFile, serverMux); err != nil {
			log.Panic("ListenAndServeTLS error", zap.Error(err))
		}
	} else {
		if err := http.ListenAndServe(s.config.Address, serverMux); err != nil {
			log.Panic("ListenAndServe error", zap.Error(err))
		}
	}
}

func (s Server) serveWs(w http.ResponseWriter, r *http.Request) {
	//requestURI := r.URL.RequestURI()
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	log.Error("ERROR", zap.Error(err))
	//	w.Header().Set("Sec-Websocket-Version", "13")
	//	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	//	return
	//}

	clientId, err := s.checkToken(r.URL.Query())
	if err != nil { //Unauthorized
		log.Info("checkToken failed, error: "+err.Error(), zap.String("clientId", clientId))
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if _, ok := s.clients.Load(clientId); ok {
		log.Debug("There is the same client connection with clientId: " + clientId)
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Websocket upgrade failed: %s"+err.Error(), zap.Error(err))
		return
	}

	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 1024),
		clientId: clientId,
		server:   s,
	}

	s.register <- client

	go client.Write()
	go client.Read()
}

func (s Server) checkToken(query url.Values) (string, error) {
	clientId := query.Get(s.config.ClientIdAlias)
	nonce := query.Get(s.config.TimeAlias)
	token := query.Get(s.config.TokenAlias)
	if len(clientId) == 0 || len(token) == 0 || len(nonce) == 0 {
		return "", fmt.Errorf("request parameters error, clientId: %s, token: %s, nonce: %s", clientId, token, nonce)
	}
	//nonceTime, err := time.Parse("2006-01-02T15:04:05.000Z")
	nonceInt, err := strconv.ParseInt(nonce, 10, 64)
	if err != nil {
		return clientId, fmt.Errorf("strconv.ParseInt: parsing %s: invalid syntax", nonce)
	}
	if now := time.Now().Unix(); nonceInt+s.config.TimeWindow < now || nonceInt-s.config.TimeWindow > now {
		return clientId, fmt.Errorf(
			"nonce is invalid, %s: %s, timeWindow: %d, now: %d, clientId: %s, token: %s",
			s.config.TimeAlias, nonce, s.config.TimeWindow, now, clientId, token,
		)
	}

	str := new(bytes.Buffer)
	//uri + nonce + node + appKey
	str.WriteString(clientId)
	str.WriteString(nonce)
	str.WriteString(cfg.Config.Node)
	str.WriteString(s.config.AppKey)
	sign, err := s.signature.Sign(str.Bytes())
	if err != nil {
		return clientId, err
	}

	//todo 测试代码
	token = sign
	if sign != token {
		return clientId, fmt.Errorf("clientId: %s, error sign, sign is %s, token: %s, nonce: %s", clientId, sign, token, nonce)
	}

	return clientId, nil
}

func (s Server) JoinCluster() error {
	//todo
	return nil
}

func (s Server) LeaveCluster() error {
	//todo
	return nil
}

func (s *Server) sendToConnection(connId string, msg string) error {
	if client, ok := s.clients.Load(connId); ok {
		select {
		case client.(*Client).send <- []byte(msg):
			//log.Info("SendToConnection " + connId + ": " + msg)
			return nil

		default:
			client.(*Client).Close()
			return errors.New("send message failed to " + connId)
		}
	}

	return errors.New("send message failed, connection: " + connId + " not found")
}

func (s Server) SendToConnections(connections []string, msg string) ([]string, error) {
	var errIds []string
	for _, connId := range connections {
		if err := s.sendToConnection(connId, msg); err != nil {
			errIds = append(errIds, connId)
		}
	}
	if len(errIds) > 0 {
		return errIds, errors.New("there are messages that failed to send to the connections")
	}

	return []string{}, nil
}

func (s Server) Broadcast(msg string) {
	s.clients.Range(func(connId, _ interface{}) bool {
		s.sendToConnection(connId.(string), msg)
		return true
	})
}

func (s Server) KickConnections(connections []string) {
	for _, connId := range connections {
		if client, ok := s.clients.Load(connId); ok {
			client.(*Client).Close()
		}
	}
}

func (s Server) KickAllConnections() {
	s.clients.Range(func(_, client interface{}) bool {
		client.(*Client).Close()
		return true
	})
}

func (s Server) CheckConnectionsOnline(connections []string) []string {
	var onlineConnections []string
	for _, connId := range connections {
		if _, ok := s.clients.Load(connId); ok {
			onlineConnections = append(onlineConnections, connId)
		}
	}

	return onlineConnections
}

func (s Server) GetAllConnections() []string {
	var connectionIds []string

	s.clients.Range(func(connId, _ interface{}) bool {
		connectionIds = append(connectionIds, connId.(string))
		return true
	})

	return connectionIds
}

func (s Server) AnyCall(method string, args json.RawMessage) (ret interface{}, err error) {
	//防止切片使用出错
	defer func() {
		if r := recover(); r != nil {
			log.Error("AnyCall panic", zap.Any("r", r))

			ret = nil
			err = errors.New("AnyCall panic")
		}
	}()

	switch method {
	default:
		return nil, errors.New("unsupported rpc method: " + method)
	}
}
