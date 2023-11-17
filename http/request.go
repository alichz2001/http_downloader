package http

import (
	"fmt"
	"strings"
)

type HTTPRequest struct {
	server string

	proto   string
	version string

	method string
	uri    string

	headers map[string]string

	body []byte
}

func (h *HTTPRequest) buildStr() string {
	//str := "HEAD /Fight.Club.1999.Bluray.720p.Farsi.Dubbed.mkv HTTP/1.1\r\nHost: localhost.com\r\nUser-Agent: curl/8.4.0\r\nAccept: */*\r\n\r\n"

	var headersStr strings.Builder

	for title, value := range h.headers {
		headersStr.WriteString(title)
		headersStr.WriteString(": ")
		headersStr.WriteString(value)
		headersStr.WriteString("\r\n")
	}

	//TODO add req header to str
	return fmt.Sprintf("%s %s %s/%s\r\nHost: %s\r\n%s\r\n", h.method, h.uri, h.proto, h.version, h.server, headersStr.String())
}
