package minio_service

import (
	"fmt"
	"os"
)

type HLSSegmentObject struct {
	StreamId string
	Name     string
	Path     string
}

func (h *HLSSegmentObject) ObjectName() string {
	return fmt.Sprintf("streams/%s/segments/%s", h.StreamId, h.Name)
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
	return fmt.Sprintf("streams/%s/thumbnails/%s", s.StreamId, "thumbnail.jpg")
}
func (s *StreamThumbnailObject) FilePath() string {
	return s.Path
}
func (s *StreamThumbnailObject) ContentType() string {
	return "image/jpeg"
}

func (s *StreamThumbnailObject) EnsurePath() bool {
	if _, err := os.Stat(s.Path); err != nil {
		return false
	}
	return true
}
