package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type HttpHandlerFunc func(w ResponseWriter, r *Request)

type HttpHandler interface {
	ServeHTTP(ResponseWriter, *Request)
}

const (
	StatusOK            int = 200
	StatusNotFound      int = 404
	StatusBadRequest    int = 400
	StatusInternalError int = 500
)

const (
	CRLF                  = "\r\n"
	DefaultReadBufferSize = 4096
)

var StatusMessage = map[int]string{
	StatusOK:            "OK",
	StatusNotFound:      "Not Found",
	StatusBadRequest:    "Bad Request",
	StatusInternalError: "Internal Server Error",
}

type Server struct {
	Addr    string
	Handler HttpHandler
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
	// Set content headers if not already set
	if r.header.Get("Content-Type") == "" {
		r.header.Set("Content-Type", "text/plain")
	}

	r.WriteHeader(status)

	if !r.wroteHeader {
		r.WriteHeaderLines()
	}

	return r.w.Write(b)
}

func (r *response) WriteHeaderLines() {
	if r.wroteHeader {
		return
	}

	s := StatusMessage[r.status]
	fmt.Fprintf(r.w, "HTTP/1.1 %d %s%s", r.status, s, CRLF)

	for k, v := range r.header {
		for _, value := range v {
			fmt.Fprintf(r.w, "%s: %s%s", k, value, CRLF)
		}
	}

	r.w.WriteString(CRLF)
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
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: a.Routes(),
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

		r := &response{conn: c, req: &req, header: make(Header), w: bufio.NewWriter(c)}

		go s.handleReq(r)
	}
}

func (s *Server) handleReq(r *response) {
	defer r.conn.Close()
	r.req.Log()
	s.Handler.ServeHTTP(r, r.req)

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
