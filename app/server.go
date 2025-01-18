package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

type HandlerFunc func(ResponseWriter, *Request)

const (
	StatusOK            int = 200
	StatusNotFound      int = 404
	StatusBadRequest    int = 400
	StatusInternalError int = 500
)

type Server struct {
	Addr    string
	Handler HandlerFunc
}

type response struct {
	conn   net.Conn
	req    *Request
	header Header
	status int

	wroteHeader bool

	w *bufio.Writer
}

type ResponseWriter interface {
	Write(status int, b []byte) (int, error)
	WriteHeader(int)
	Header() Header
}

func (r *response) Header() Header {
	return r.header
}

func (r *response) Write(status int, b []byte) (int, error) {
	r.WriteHeader(status)
	r.Header().Set("Content-Type", "text/plain")
	r.Header().Set("Content-Length", strconv.Itoa(len(b)))

	if !r.wroteHeader {
		r.WriteHeaderLines()
	}
	return r.w.Write(b)
}

func (r *response) WriteHeaderLines() {
	if r.wroteHeader {
		return
	}

	var s string
	switch r.status {
	case StatusOK:
		s = "OK"
	case StatusNotFound:
		s = "Not Found"
	case StatusBadRequest:
		s = "Bad Request"
	case StatusInternalError:
		s = "Internal Server Error"
	}
	fmt.Fprintf(r.w, "HTTP/1.1 %d %s\r\n", r.status, s)

	for k, v := range r.header {
		for _, value := range v {
			fmt.Fprintf(r.w, "%s: %s\r\n", k, value)
		}
	}

	r.w.WriteString("\r\n")
	r.wroteHeader = true
}

func (w *response) WriteHeader(code int) {
	checkWriteHeader(code)
	w.status = code
}

func (r *response) Flush() error {
	return r.w.Flush()
}

func (a *application) Serve() error {
	port := a.config.Port
	s := &Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
	}

	fmt.Printf("Server started on address: %s\n", s.Addr)

	return s.ListenAndServe()
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind to address: %s", s.Addr)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		req, err := ParseRequest(c)

		if err != nil {
			fmt.Println("Error reading headers: ", err.Error())
			// TODO: Handle request parsing error
			continue
		}

		rw := &response{conn: c, req: &req, header: make(Header), w: bufio.NewWriter(c)}

		go handleReq(rw)
	}
}

func handleReq(r *response) {
	defer r.conn.Close()
	r.req.Log()
	HandleRoute(r, r.req)

	if err := r.Flush(); err != nil {
		fmt.Println("Error flushing response:", err)
	}
}

func (req *Request) Log() {
	// Log request in framework style
	timestamp := time.Now().Format("2025/01/10 15:04:05")
	userAgent := "N/A"
	host := "N/A"
	contentType := "N/A"

	if ua, ok := req.Header["User-Agent"]; ok {
		userAgent = fmt.Sprintf("%v", ua)
	}

	if h, ok := req.Header["Host"]; ok {
		host = fmt.Sprintf("%v", h)
	}

	if ct, ok := req.Header["Content-Type"]; ok {
		contentType = fmt.Sprintf("%v", ct)
	}

	fmt.Printf("[%s] %s %s %s | Host: %s | User-Agent: %s | Content-Type: %s\n",
		timestamp, req.Method, req.Path, req.Proto, host, userAgent, contentType)
}

func checkWriteHeader(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
