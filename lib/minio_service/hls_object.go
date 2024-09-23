package minio_service

import "fmt"

type HLSSegmentObject struct {
	StreamId string
	Name     string
	Path     string
}

func (h *HLSSegmentObject) ObjectName() string {
	return fmt.Sprintf("streams/%s/%s", h.StreamId, h.Name)
}

func (h *HLSSegmentObject) FilePath() string {
	return h.Path
}

func (h *HLSSegmentObject) ContentType() string {
	return "video/mp2t"
}

type StreamThumbnailObject struct {
	StreamId string
	Path     string
}

func (s *StreamThumbnailObject) ObjectName() string {
	return fmt.Sprintf("thumbnails/%s/%s", s.StreamId, "thumbnail.jpg")
}
func (s *StreamThumbnailObject) FilePath() string {
	return s.Path
}
func (s *StreamThumbnailObject) ContentType() string {
	return "image/jpeg"
}
