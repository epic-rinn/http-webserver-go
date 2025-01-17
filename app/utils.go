package main

import (
	"fmt"
	"net"
	"strings"
)

type Headers = map[string]any

type Request struct {
	Method  string
	Path    string
	Version string
	Headers Headers
}

func ParseRequest(c net.Conn) (Request, error) {
	tempBuf := make([]byte, 4096)
	var cd strings.Builder

	n, err := c.Read(tempBuf)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return Request{}, err
	}

	req := Request{}
	headers := make(Headers)

	cd.Write(tempBuf[:n])
	lines := strings.Split(cd.String(), "\r\n")

	for i, line := range lines {
		if line == "" {
			break
		}

		if i == 0 {
			parts := strings.SplitN(line, " ", 3)
			req.Method = parts[0]
			req.Path = parts[1]
			req.Version = parts[2]
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 2 {
			headers[parts[0]] = strings.TrimSpace(parts[1])
		}
	}
	req.Headers = headers

	return req, nil
}
