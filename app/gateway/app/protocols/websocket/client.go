package websocket

import (
	"time"

	"github.com/gopusher/gateway/pkg/log"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

//Client is websocket client
type Client struct {
	clientId string
	conn     *websocket.Conn
	send     chan []byte
	server   Server
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

//var (
//	newline = []byte{'\n'}
//)

//Write write to client from message channel
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			log.Error("conn Close err: "+err.Error(), zap.Error(err))
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Error("SetWriteDeadline err: "+err.Error(), zap.Error(err))
			}
			if !ok {
				// The hub closed the channel.
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Error("WriteMessage err: "+err.Error(), zap.Error(err))
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error("NextWriter err: "+err.Error(), zap.Error(err))
				return
			}
			if _, err := w.Write(message); err != nil {
				log.Error("Write err: "+err.Error(), zap.Error(err))
			}

			// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

			if err := w.Close(); err != nil {
				log.Error("Close err: "+err.Error(), zap.Error(err))
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Error("SetWriteDeadline err: "+err.Error(), zap.Error(err))
				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Error("WriteMessage err: "+err.Error(), zap.Error(err))
				return
			}
		}
	}
}

//Read from client
func (c *Client) Read() {
	defer func() {
		c.server.unregister <- c
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Error("SetReadDeadline err: "+err.Error(), zap.Error(err))
		return
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Error("SetReadDeadline err: "+err.Error(), zap.Error(err))
		}
		return nil
	})
	for {
		//_, message, err := c.conn.ReadMessage()
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Error("IsUnexpectedCloseError err: "+err.Error(), zap.Error(err))
			}
			break
		}
		//log.Info("msg :" + string(message))
	}
}

//Close client connection
func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		log.Error("Close err: "+err.Error(), zap.Error(err))
	}
	return
}

//SendMessage send message to client
func (c *Client) SendMessage(message string) bool {
	return true
}
