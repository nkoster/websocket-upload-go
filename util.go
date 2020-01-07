package main

import (
	"net/http"
	"os"
)

func mimeType(filename string) string {
	if f, err := os.Open(filename); err == nil {
		buffer := make([]byte, 512)
		if _, err := f.Read(buffer); err == nil {
			return http.DetectContentType(buffer)
		}
	}
	return ""
}
