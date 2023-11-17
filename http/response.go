package http

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type HTTPResponse struct {
	proto   string
	version string

	// this is not for HTTP standard. for read response body should know request method.
	method string

	statusCode         string
	statusReasonPhrase string

	headers map[string]string
}

func NewHTTPResponse() *HTTPResponse {
	return &HTTPResponse{
		headers: make(map[string]string),
	}
}

func (h *HTTPResponse) readResponse(reader io.Reader) {
	r := bufio.NewReader(reader)

	statusLine, _ := r.ReadString('\r')
	statusLine = statusLine[:len(statusLine)-1]
	//log.Printf(statusLine)
	h.parseStatusLine(statusLine)
	r.ReadByte()
	var headerLine string
	for {
		headerLine, _ = r.ReadString('\r')
		//log.Printf(headerLine)
		r.ReadByte()
		if headerLine != "\r" {
			headerLine = headerLine[:len(headerLine)-1]
			splitedHeader := strings.Split(headerLine, ": ")
			h.headers[strings.ToLower(splitedHeader[0])] = splitedHeader[1]
		} else {
			// do not read body of request if method is HEAD
			// TODO maybe should not read body for some other methods
			if h.method == "HEAD" {
				break
			}

			break
		}
	}
}

func (h *HTTPResponse) parseStatusLine(statusLine string) {
	var tmp []string
	tmp = strings.Split(statusLine, "/")
	h.proto = tmp[0]

	tmp = strings.SplitN(tmp[1], " ", 3)
	h.version = tmp[0]
	h.statusCode = tmp[1]
	h.statusReasonPhrase = tmp[2]
}

func (h *HTTPResponse) GetHeader(headerStr string) (string, error) {
	res, ok := h.headers[headerStr]
	if ok {
		return res, nil
	}
	return "", fmt.Errorf("not found")
}

func (h *HTTPResponse) GetStatusCode() string {
	return h.statusCode
}
