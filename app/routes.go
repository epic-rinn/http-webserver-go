package main

import "net"

func HandleRoute(c net.Conn, headers Request) {
	if headers["Path"] == "/" {
		c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
