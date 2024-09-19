package main

import "fmt"

func BuildHLSUrl(stream_id string, m3u8_file string) string {
	cfg := GetConfig()
	return fmt.Sprintf("%s/%s/%s", cfg.HLS_PUBLIC_URL, stream_id, m3u8_file)
}
