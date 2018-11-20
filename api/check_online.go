package api

func (s *Server) CheckConnectionsOnline(message *ConnectionsMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	var onlineConnections []string
	if len(message.Connections) > 0 {
		onlineConnections = s.server.CheckConnectionsOnline(message.Connections)
	}

	*reply = s.success(onlineConnections)
	return nil
}

func (s *Server) GetAllConnections(message *TokenMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	onlineConnections := s.server.GetAllConnections()

	*reply = s.success(onlineConnections)
	return nil
}
