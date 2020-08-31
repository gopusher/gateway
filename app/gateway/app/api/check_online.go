package api

//CheckConnectionsOnline check connections is online
func (s *Server) CheckConnectionsOnline(message *ConnectionsMessage, reply *Response) error {
	if errResp := s.checkToken(message.Token); errResp != nil {
		*reply = *errResp
		return nil
	}

	var onlineConnections []string
	if len(message.Connections) > 0 {
		onlineConnections = s.server.CheckConnectionsOnline(message.Connections)
	}

	*reply = s.success(onlineConnections)
	return nil
}

//GetAllConnections returns all connections
func (s *Server) GetAllConnections(message *TokenMessage, reply *Response) error {
	if errResp := s.checkToken(message.Token); errResp != nil {
		*reply = *errResp
		return nil
	}

	onlineConnections := s.server.GetAllConnections()

	*reply = s.success(onlineConnections)
	return nil
}
