package main

import (
	"fmt"
	"strings"
)

var httpStatusText = map[int]string{
	200: "OK",
	204: "No Content",
	404: "Not Found",
	405: "Method Not Allowed",
	500: "Internal Server Error",
}

type response struct {
	body    []byte
	headers map[string]string
	status  int
}

func (r *response) getBytes() []byte {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.status, httpStatusText[r.status]))
	for key, value := range r.headers {
		b.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	b.WriteString("\r\n")
	return append([]byte(b.String()), r.body...)
}

func responseWithStatus(status int) response {
	return response{
		body:    make([]byte, 0),
		headers: make(map[string]string),
		status:  status,
	}
}
