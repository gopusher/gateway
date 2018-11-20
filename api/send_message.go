package api

type Message struct {
	ConnectionsMessage
	Msg 			string		`json:"msg"` 		//为一个json，里边包含 type 消息类型
}

func (s *Server) SendToConnections(message *Message, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	if len(message.Connections) == 0 {
		*reply = s.failure(nil, "empty connections")
		return nil
	}

	if message.Msg == "" {
		*reply = s.failure(nil, "empty msg")
		return nil
	}

	if errIds, err := s.server.SendToConnections(message.Connections, message.Msg); err != nil {
		*reply = s.failure(errIds, err.Error())
		return nil
	}

	*reply = s.success(nil)
	return nil
}

type BroadcastMessage struct {
	TokenMessage
	Msg 			string		`json:"msg"` 		//为一个json，里边包含 type 消息类型
}

func (s *Server) Broadcast(message *BroadcastMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	if message.Msg == "" {
		*reply = s.failure(nil, "empty msg")
		return nil
	}

	go s.server.Broadcast(message.Msg)

	*reply = s.success(nil)
	return nil
}
