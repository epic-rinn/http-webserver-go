package main

import (
	"fmt"
	"net"
	"net/textproto"
	"strings"
)

type Header map[string][]string

func (h Header) Get(key string) string {
	if v, ok := h[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
// The key is case insensitive; it is canonicalized by
// [CanonicalHeaderKey].
func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}

// Set sets the header entries associated with key to the
// single element value. It replaces any existing values
// associated with key. The key is case insensitive; it is
// canonicalized by [textproto.CanonicalMIMEHeaderKey].
// To use non-canonical keys, assign to the map directly.
func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

// Del deletes the values associated with key.
// The key is case insensitive; it is canonicalized by
// [CanonicalHeaderKey].
func (h Header) Del(key string) {
	textproto.MIMEHeader(h).Del(key)
}

type Request struct {
	// Method specifies the HTTP method (GET, POST, PUT, etc.).
	// For client requests, an empty string means GET.
	Method string

	// Path specifies the path (relative paths may omit leading slash)
	Path string

	// The protocol version for incoming server requests.
	Proto string

	// Header contains the request header fields either received
	// by the server or to be sent by the client.
	// If a server received a request with header lines,
	//
	//	Host: example.com
	//	accept-encoding: gzip, deflate
	//	Accept-Language: en-us
	//	fOO: Bar
	//	foo: two
	//
	// then
	//
	//	Header = map[string][]string{
	//		"Host":           {"example.com"},
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Foo": {"two"},
	//	}
	Header Header
}

func ParseRequest(c net.Conn) (Request, error) {
	tempBuf := make([]byte, 4096)
	var cd strings.Builder

	n, err := c.Read(tempBuf)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return Request{}, err
	}

	req := Request{}
	header := make(Header)

	cd.Write(tempBuf[:n])
	lines := strings.Split(cd.String(), "\r\n")

	for i, line := range lines {
		if line == "" {
			break
		}

		if i == 0 {
			parts := strings.SplitN(line, " ", 3)
			req.Method = parts[0]
			req.Path = parts[1]
			req.Proto = parts[2]
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 2 {
			prev := header[parts[0]]
			header[parts[0]] = append(prev, strings.TrimSpace(parts[1]))
		}
	}
	req.Header = header

	return req, nil
}
