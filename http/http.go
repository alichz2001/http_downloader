package http

import (
	"fmt"
	"io"
	"log"
	"net"
)

type HTTP struct {
	server string

	request  *HTTPRequest
	response *HTTPResponse

	conn *net.TCPConn
}

func NewHTTP(server string, uri string, method string) *HTTP {
	return &HTTP{
		server: server,
		request: &HTTPRequest{
			server:  server,
			proto:   "HTTP",
			version: "1.1",
			method:  method,
			uri:     uri,
			headers: make(map[string]string),
		},
	}
}

func (h *HTTP) Run() {
	//server, _ := net.ResolveTCPAddr("tcp", h.server)
	//client, _ := net.ResolveTCPAddr("tcp", ":50013")
	//conn, err := net.DialTCP("tcp", client, server)
	conn, err := net.Dial("tcp", h.server)
	if err != nil {
		log.Fatalf("%s", err)
	}

	h.conn = conn.(*net.TCPConn)
	log.Printf("%#v", h.request)

	fmt.Fprintf(conn, h.request.buildStr())
	h.parseResponse()
}

func (h *HTTP) Close() error {
	return h.conn.Close()
}

func (h *HTTP) parseResponse() {
	h.response = &HTTPResponse{
		method:  h.request.method,
		headers: make(map[string]string),
	}
	h.response.readResponse(h.conn)
}

func (h *HTTP) GetResponse() *HTTPResponse {
	return h.response
}

func (h *HTTP) GetRequestMethod() string {
	return h.request.method
}

func (h *HTTP) GetResponseBodyWriter() io.Reader {
	return h.conn
}

func (h *HTTP) SetRequestHeader(title string, value string) *HTTP {
	h.request.headers[title] = value
	return h
}
