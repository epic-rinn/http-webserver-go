package main

import (
	"fmt"
	"net"
	"time"
)

func (a *application) Serve() error {
	port := a.config.Port
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return fmt.Errorf("failed to bind to port: %d", port)
	}

	fmt.Printf("Server started on port :%d\n", port)

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		go handleReq(c)
	}
}

func handleReq(c net.Conn) {
	req, err := ParseRequest(c)
	if err != nil {
		fmt.Println("Error reading headers: ", err.Error())
		return
	}

	logReq(req)

	HandleRoute(c, req)
	c.Close()
}

func logReq(req Request) {
	// Log request in framework style
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	userAgent := "N/A"
	host := "N/A"
	contentType := "N/A"

	if ua, ok := req.Headers["User-Agent"]; ok {
		userAgent = fmt.Sprintf("%v", ua)
	}

	if h, ok := req.Headers["Host"]; ok {
		host = fmt.Sprintf("%v", h)
	}

	if ct, ok := req.Headers["Content-Type"]; ok {
		contentType = fmt.Sprintf("%v", ct)
	}

	fmt.Printf("[%s] %s %s %s | Host: %s | User-Agent: %s | Content-Type: %s\n",
		timestamp, req.Method, req.Path, req.Version, host, userAgent, contentType)
}
