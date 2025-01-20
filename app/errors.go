package main

import (
	"bufio"
	"fmt"
	"net"
)

func (app *application) NotFound(w ResponseWriter, r *Request) {
	w.Write(StatusNotFound, []byte("Not Found"))
}

func (s *Server) HandleUnknownError(c net.Conn, status int, message string) {
	w := bufio.NewWriter(c)
	statusMsg := StatusMessage[status]

	fmt.Fprintf(w, "HTTP/1.1 %d %s%s", status, statusMsg, CRLF)
	fmt.Fprintf(w, "Content-Type: text/plain%s", CRLF)
	fmt.Fprintf(w, "Content-Length: %d%s", len(message), CRLF)
	fmt.Fprintf(w, "Connection: close%s", CRLF)
	w.WriteString(CRLF)
	w.WriteString(message)
	w.Flush()
}
