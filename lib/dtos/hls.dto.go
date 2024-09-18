package dtos

type HLSPublishStartInput struct {
	StreamId        string `json:"stream_id"`
	StreamServer    string `json:"stream_server"`
	StreamServerApp string `json:"stream_server_app"`
}

type HLSPublishStartResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type HLSPublishStopInput struct {
	StreamId string `json:"stream_id"`
}
