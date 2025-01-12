package main

import "net"

func HandleRoute(c net.Conn, headers Request) {
	switch headers["Path"] {
	case "/":
		c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
