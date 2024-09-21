package main

import (
	"fmt"
	"strings"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/utils"
)

const PUBLIC_DIR = "public"
const M3U8_FILE = "playlist.m3u8"
const SEGMENT_FILE_PREFIX = "seg-"
const SEGMENT_FILE_SUFFIX = ".ts"
const SEGMENT_FILE = SEGMENT_FILE_PREFIX + "%05d" + SEGMENT_FILE_SUFFIX
const THUMBNAIL_DIR = "thumbnail"
const THUMBNAIL_FILE = "thumbnail.jpg"

func BuildHLSUrl(stream_id string) string {
	cfg := configs.GetConfig()
	return fmt.Sprintf("%s/%s/%s", cfg.HLS_PUBLIC_URL(), stream_id, M3U8_FILE)
}

func BuildHLSStreamDir(stream_id string) string {
	return fmt.Sprintf("%s/%s", PUBLIC_DIR, stream_id)
}

func BuildHLSm3u8FilePath(stream_id string) string {
	return fmt.Sprintf("%s/%s/%s", PUBLIC_DIR, stream_id, M3U8_FILE)
}

func BuildHLSSegmentFilePath(stream_id string) string {
	return fmt.Sprintf("%s/%s/%s", PUBLIC_DIR, stream_id, SEGMENT_FILE)
}

func BuildStreamURL(server_url string, server_app string, stream_id string, subscribe_key string) string {
	return fmt.Sprintf("%s?streamid=%s%s?key=%s", server_url, server_app, stream_id, subscribe_key)
}

func BuildThumbnailFilePath(stream_id string) string {
	return fmt.Sprintf("%s/%s/%s/%s", PUBLIC_DIR, stream_id, THUMBNAIL_DIR, THUMBNAIL_FILE)
}

func BuildThumbnailUrl(stream_id string) string {
	cfg := configs.GetConfig()
	return fmt.Sprintf("%s/%s/%s/%s", cfg.HLS_PUBLIC_URL(), stream_id, THUMBNAIL_DIR, THUMBNAIL_FILE)
}

func GetListSegmentFilesByStreamId(stream_id string) *[]string {
	files, err := utils.ListDirWithFilter(BuildHLSStreamDir(stream_id), func(f_name string) bool {
		return strings.HasPrefix(f_name, SEGMENT_FILE_PREFIX) && strings.HasSuffix(f_name, SEGMENT_FILE_SUFFIX)
	})
	if err != nil {
		return &[]string{}
	}
	return files
}
