package api

//KickConnections kick connections
func (s *Server) KickConnections(message *ConnectionsMessage, reply *Response) error {
	if errResp := s.checkToken(message.Token); errResp != nil {
		*reply = *errResp
		return nil
	}

	go s.server.KickConnections(message.Connections)

	*reply = s.success(nil)
	return nil
}

//KickAllConnections kick all connections
func (s *Server) KickAllConnections(message *TokenMessage, reply *Response) error {
	if errResp := s.checkToken(message.Token); errResp != nil {
		*reply = *errResp
		return nil
	}

	go s.server.KickAllConnections()

	*reply = s.success(nil)
	return nil
}
