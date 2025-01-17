package main

import (
	"fmt"
	"net"
)

func (a *application) Serve() error {
	port := a.config.Port
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return fmt.Errorf("failed to bind to port: %d", port)
	}

	fmt.Printf("Server started on port :%d", port)

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
	headers, err := GetHeaders(c)
	if err != nil {
		fmt.Println("Error reading headers: ", err.Error())
		return
	}

	fmt.Println(headers)

	HandleRoute(c, headers)
	c.Close()
}
