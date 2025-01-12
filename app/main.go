package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var _ = net.Listen
var _ = os.Exit
var wg sync.WaitGroup

func handleConn(c net.Conn) {
	tempBuf := make([]byte, 4096)
	var cd strings.Builder

	n, err := c.Read(tempBuf)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	cd.Write(tempBuf[:n])
	fmt.Println(cd.String())

	c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	c.Close()

	wg.Done()
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	fmt.Println("Server started on port :4221")

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		wg.Add(1)
		go handleConn(c)
		wg.Wait()
	}
}
