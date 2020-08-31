package api

import (
	"encoding/json"
)

type AnyCallMessage struct {
	TokenMessage
	Method string          `json:"method"`
	Args   json.RawMessage `json:"args"`
}

//获取部分 order book
func (s *Server) AnyCall(message *AnyCallMessage, reply *Response) error {
	if errResp := s.checkToken(message.Token); errResp != nil {
		*reply = *errResp
		return nil
	}
	//log.Debug("AnyCall method: " + message.Method + ", args: " + string(message.Args))

	data, err := s.server.AnyCall(message.Method, message.Args)
	if err != nil {
		*reply = s.failure(ServerErrorCode, err.Error())
		return nil
	}

	*reply = s.success(data)
	return nil
}
