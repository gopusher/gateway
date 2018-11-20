package api

func (s *Server) KickConnections(message *ConnectionsMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	go s.server.KickConnections(message.Connections)

	*reply = s.success(nil)
	return nil
}

func (s *Server) KickAllConnections(message *TokenMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	go s.server.KickAllConnections()

	*reply = s.success(nil)
	return nil
}
