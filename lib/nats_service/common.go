package nats_service

import (
	"encoding/json"
	"time"
)

const DEFAULT_TIMEOUT time.Duration = time.Second * 20

type IRequestMessage interface {
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
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func (m *ResponseMessage) Encode() ([]byte, error) {
	return json.Marshal(m)
}
