package api

func (s *Server) GetNodeId(message *TokenMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	*reply = s.success([]string{s.nodeId})
	return nil
}
