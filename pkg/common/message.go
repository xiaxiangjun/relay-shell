package common

import "encoding/json"

const MsgResponse = "response"
const MsgConnect = "connect"
const MsgRegister = "register"

type Message struct {
	Type    string          `json:"type"` // register, connect
	User    string          `json:"user,omitempty"`
	Time    int64           `json:"ts,omitempty"`
	UUID    string          `json:"uuid,omitempty"` // 每次请求都不能一样，服务器在这里做校验
	Sign    string          `json:"sign,omitempty"` // hmac-sha1(user+":"+ts+":"+type+":"+uuid, password)
	Code    int             `json:"code,omitempty"`
	Msg     string          `json:"msg,omitempty"`
	Payload json.RawMessage `json:"payload"`
}

type RegisterRequest struct {
	SessionID string `json:"sid"`
}

type RegisterResponse struct{}

type ConnectRequest struct {
	TargetSID  string `json:"sid"`
	Mode       string `json:"mode"` // ssh, tunnel
	TunnelAddr string `json:"addr"`
	Sign       string `json:"sign,omitempty"` // hmac-sha1(mode+":"+sid+":"+addr, password)
}

type ConnectResponse struct {
}

func MarshalToJson(msg *Message, payload interface{}) ([]byte, error) {
	msg.Payload, _ = json.Marshal(payload)
	return json.Marshal(msg)
}
