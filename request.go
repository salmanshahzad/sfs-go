package main

import (
	"bytes"
	"errors"
)

type request struct {
	body     []byte
	headers  map[string]string
	method   string
	resource string
}

func parseRequest(buf []byte) (request, error) {
	req := request{}
	method, resource, remaining, err := parseMethodAndResource(buf)
	if err != nil {
		return req, err
	}

	headers, remaining, err := parseHeaders(remaining)
	if err != nil {
		return req, err
	}

	req.body = remaining
	req.headers = headers
	req.method = method
	req.resource = resource
	return req, nil
}

func parseMethodAndResource(buf []byte) (method, resource string, remaining []byte, err error) {
	before, after, found := bytes.Cut(buf, []byte("\r\n"))
	if !found {
		err = errors.New("Invalid request")
		return
	}

	parts := bytes.Split(before, []byte(" "))
	if len(parts) != 3 {
		err = errors.New("Invalid request")
		return
	}

	if string(parts[2]) != "HTTP/1.1" {
		err = errors.New("Invalid request")
		return
	}

	method = string(parts[0])
	resource = string(parts[1])
	remaining = after
	return
}

func parseHeaders(buf []byte) (headers map[string]string, remaining []byte, err error) {
	headers = make(map[string]string)
	remaining = buf
	for {
		before, after, found := bytes.Cut(remaining, []byte("\r\n"))
		if !found {
			err = errors.New("Invalid headers")
			return
		}

		remaining = after
		if len(before) == 0 {
			return
		}

		parts := bytes.Split(before, []byte(": "))
		if len(parts) != 2 {
			err = errors.New("Invalid headers")
			return
		}

		headers[string(parts[0])] = string(parts[1])
	}
}
