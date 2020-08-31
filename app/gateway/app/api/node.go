package api

//GetNodeId get the comet node id
func (s *Server) GetNode(message *TokenMessage, reply *Response) error {
	if errResp := s.checkToken(message.Token); errResp != nil {
		*reply = *errResp
		return nil
	}

	*reply = s.success(s.node)
	return nil
}
