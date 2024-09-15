package nats_service

import (
	"encoding/json"
	"time"
)

const DEFAULT_TIMEOUT time.Duration = time.Second * 20

type IMessage interface {
	Data() []byte
	JSONParse(interface{}) error
}

type RequestMessage struct {
	data []byte
}

func (m *RequestMessage) Data() []byte {
	return m.data
}

func (m *RequestMessage) JSONParse(out interface{}) error {
	return json.Unmarshal(m.data, out)
}

type ResponseMessage struct {
	Data  []byte `json:"data"`
	Error error  `json:"error"`
}

func (m *ResponseMessage) Encode() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ResponseMessage) Decode(out interface{}) error {
	return json.Unmarshal(m.Data, out)
}
