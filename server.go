package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

func handleRequest(req *request, dir string) response {
	if strings.HasSuffix(req.resource, "/") {
		req.resource = path.Join(req.resource, "index.html")
		return handleRequest(req, dir)
	}

	file := strings.TrimLeft(req.resource, "/")
	file = path.Join(dir, file)

	stat, err := os.Stat(file)
	if err != nil {
		return responseWithStatus(http.StatusNotFound)
	}

	if req.method != http.MethodHead && req.method != http.MethodGet {
		return responseWithStatus(http.StatusMethodNotAllowed)
	}

	res := response{
		body: make([]byte, 0),
		headers: map[string]string{
			"Content-Length": fmt.Sprint(stat.Size()),
			"Content-Type":   getContentType(file),
		},
		status: http.StatusOK,
	}
	if req.method == http.MethodGet {
		buf, err := os.ReadFile(file)
		if err != nil {
			stderrPrint("Error reading file", file, err)
			return responseWithStatus(http.StatusInternalServerError)
		}
		res.body = buf
	}

	return res
}

func getContentType(file string) string {
	switch path.Ext(file) {
	case ".css":
		return "text/css"
	case ".gif":
		return "image/gif"
	case ".html":
		return "text/html"
	case ".jpg":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".js":
		return "text/javascript"
	case ".md":
		return "text/markdown"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
