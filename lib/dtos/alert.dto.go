package dtos

type AlertPayload struct {
	IP      string
	Method  string
	Url     string
	Payload string
}

type StreamStopAlert struct {
	StreamId string
	Title    string
}
