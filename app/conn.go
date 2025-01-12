package main

import (
	"fmt"
	"net"
)

func HandleConn(c net.Conn) {
	headers, err := GetHeaders(c)
	if err != nil {
		fmt.Println("Error reading headers: ", err.Error())
		return
	}

	fmt.Println(headers)

	HandleRoute(c, headers)
	c.Close()
}
