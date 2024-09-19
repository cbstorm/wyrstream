package main

import "fmt"

const PUBLIC_DIR = "public"
const M3U8_FILE = "playlist.m3u8"
const SEGMENT_FILE = "seg-%05d.ts"
const SEGMENT_FILE_PREFIX = "seg-"
const SEGMENT_FILE_SUFFIX = ".ts"

func BuildHLSUrl(stream_id string) string {
	cfg := GetConfig()
	return fmt.Sprintf("%s/%s/%s", cfg.HLS_PUBLIC_URL, stream_id, M3U8_FILE)
}

func BuildHLSStreamDir(stream_id string) string {
	return fmt.Sprintf("%s/%s", PUBLIC_DIR, stream_id)
}

func BuildHLSm3u8FilePath(stream_id string) string {
	return fmt.Sprintf("%s/%s/%s", PUBLIC_DIR, stream_id, M3U8_FILE)
}

func BuildHLSSegmentFile(stream_id string) string {
	return fmt.Sprintf("%s/%s/%s", PUBLIC_DIR, stream_id, SEGMENT_FILE)
}

func BuildStreamURL(server_url string, server_app string, stream_id string, subscribe_key string) string {
	return fmt.Sprintf("%s?streamid=%s%s?key=%s", server_url, server_app, stream_id, subscribe_key)
}
