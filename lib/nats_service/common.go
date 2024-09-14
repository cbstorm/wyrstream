package nats_service

import "encoding/json"

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
