package main

import (
	"fmt"
	"net"
	"strings"
)

type Request = map[string]string

func GetHeaders(c net.Conn) (Request, error) {
	tempBuf := make([]byte, 4096)
	var cd strings.Builder

	n, err := c.Read(tempBuf)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return nil, err
	}

	headers := make(map[string]string)

	cd.Write(tempBuf[:n])
	lines := strings.Split(cd.String(), "\r\n")

	for i, line := range lines {
		if i == 0 {
			parts := strings.SplitN(line, " ", 3)
			headers["Method"] = parts[0]
			headers["Path"] = parts[1]
			headers["Version"] = parts[2]
		}
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 2 {
			headers[parts[0]] = strings.TrimSpace(parts[1])
		}
	}

	return headers, nil
}
