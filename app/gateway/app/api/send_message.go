package api

//Message is api message
type Message struct {
	ConnectionsMessage
	Msg string `json:"msg"` //为一个json，里边包含 type 消息类型
}

//SendToConnections send message to connections
func (s *Server) SendToConnections(message *Message, reply *Response) error {
	if errResp := s.checkToken(message.Token); errResp != nil {
		*reply = *errResp
		return nil
	}

	if len(message.Connections) == 0 {
		*reply = s.failure(ParametersError, "empty connections")
		return nil
	}

	if message.Msg == "" {
		*reply = s.failure(ParametersError, "empty msg")
		return nil
	}

	if errIds, err := s.server.SendToConnections(message.Connections, message.Msg); err != nil {
		*reply = s.failureWithData(SendToConnectionsErrorCode, errIds, err.Error())
		return nil
	}

	*reply = s.success(nil)
	return nil
}

//BroadcastMessage is a broadcast type message
type BroadcastMessage struct {
	TokenMessage
	Msg string `json:"msg"` //为一个json，里边包含 type 消息类型
}

//Broadcast send message to all connections
func (s *Server) Broadcast(message *BroadcastMessage, reply *Response) error {
	if errResp := s.checkToken(message.Token); errResp != nil {
		*reply = *errResp
		return nil
	}

	if message.Msg == "" {
		*reply = s.failure(ParametersError, "empty msg")
		return nil
	}

	go s.server.Broadcast(message.Msg)

	*reply = s.success(nil)
	return nil
}
