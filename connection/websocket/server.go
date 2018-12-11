package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/gopusher/gateway/configuration"
	"errors"
	"github.com/gopusher/gateway/notification"
	"github.com/gopusher/gateway/log"
	"time"
)

type Server struct {
	config *configuration.CometConfig
	rpc *notification.Client
	upgrader websocket.Upgrader
	register chan *Client
	unregister chan *Client
	clients map[string]*Client
}

func NewWebSocketServer(config *configuration.CometConfig) *Server {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	rpc := notification.NewRpc(config.NotificationUrl, config.NotificationUserAgent)
	return &Server{
		config: config,
		rpc: rpc,
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

// run websocket server
func (s *Server) initWsServer() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ws", s.serveWs)

	log.Info("Websocket server start running with SocketProtocol: %s, Listen:\"%s\"", s.config.SocketProtocol, s.config.SocketPort)
	if s.config.SocketProtocol == "wss" {
		if err := http.ListenAndServeTLS(s.config.SocketPort, s.config.SocketCertFile, s.config.SocketKeyFile, serverMux); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(s.config.SocketPort, serverMux); err != nil {
			panic(err)
		}
	}
}

func (s *Server) handleClients() {
	for {
		select {
		case client := <-s.register:
			s.clients[client.connId] = client

			//notify router api server
			if _, err := s.rpc.Call("Online", client.connId, s.config.NodeId); err != nil {
				log.Error("Client online notification failed: %s", err.Error())
			}
		case client := <-s.unregister:
			if _, ok := s.clients[client.connId]; ok {
				delete(s.clients, client.connId)
				close(client.send)
				client.Close()

				//notify router api server
				if _, err := s.rpc.Call("Offline", client.connId, s.config.NodeId); err != nil {
					log.Error("Client online notification failed: %s", err.Error())
				}
			}
		}
	}
}

func (s Server) serveWs(w http.ResponseWriter, r *http.Request) {
	connId, err := s.checkToken(r.URL.Query())
	if err != nil {
		log.Error("Check token failed, connId: %s, error: %s", connId, err.Error())

		//Unauthorized
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	//There is the same connId client connection
	if _, ok := s.clients[connId]; ok {
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Websocket upgrade failed: %s", err.Error())
		return
	}

	client := &Client{
		conn: c,
		send: make(chan []byte, 1024),
		connId: connId,
		server: s,
	}

	s.register <- client

	go client.Write()
	go client.Read()
}

func (s Server) checkToken(query map[string][]string) (string, error) {
	if c, ok := query["c"]; !ok || len(c) < 1 {
		return "", errors.New("need param c")
	}
	connId := query["c"][0]

	if t, ok := query["t"]; !ok || len(t) < 1 {
		return connId, errors.New("need param t")
	}
	token := query["t"][0]

	if _, err := s.rpc.Call("CheckToken", connId, token, s.config.NodeId); err != nil {
		return "", err
	}

	return connId, nil
}

func (s *Server) SendToConnection(connId string, msg string) error {
	if client, ok := s.clients[connId]; ok {
		select {
		case client.send <- []byte(msg):
			//log.Info("SendToConnection " + connId + ": " + msg)
			return nil

		default:
			client.Close()
			return errors.New("send message failed to " + connId)
		}
	}

	return errors.New("send message failed, connection: " + connId + " not found")
}

func (s *Server) SendToConnections(connections []string, msg string) ([]string, error) {
	var errIds []string
	for _, connId := range connections {
		if err := s.SendToConnection(connId, msg); err != nil {
			errIds = append(errIds, connId)
		}
	}
	if len(errIds) > 0 {
		return errIds, errors.New("there are messages that failed to send to the connections")
	}

	return []string{}, nil
}

func (s *Server) Broadcast(msg string) {
	for connId := range s.clients {
		s.SendToConnection(connId, msg)
	}
}

func (s *Server) KickConnections(connections []string) {
	for _, connId := range connections {
		if client, ok := s.clients[connId]; ok {
			client.Close()
		}
	}
}

func (s *Server) KickAllConnections() {
	for _, client := range s.clients {
		client.Close()
	}
}

func (s *Server) CheckConnectionsOnline(connections []string) []string {
	var onlineConnections []string
	for _, connId := range connections {
		if _, ok := s.clients[connId]; ok {
			onlineConnections = append(onlineConnections, connId)
		}
	}

	return onlineConnections
}

func (s *Server) GetAllConnections() []string {
	var connectionIds []string

	for connId := range s.clients {
		connectionIds = append(connectionIds, connId)
	}

	return connectionIds
}

func (s *Server) JoinCluster() {
	//wait for rpc and ws server bootstrap
	time.Sleep(time.Duration(5)*time.Second)
	//notify router api server
	if _, err := s.rpc.Call("JoinCluster", s.config.NodeId); err != nil {
		log.Error("Gateway JoinCluster notification failed: %s", err.Error())
		return
	}

	log.Info("Gateway JoinCluster notification success")
}
