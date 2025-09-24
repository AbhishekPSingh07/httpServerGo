package request

import (
	"fmt"
	"io"
	"strings"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) ValidMethod() bool {
	return r.Method == "GET" || r.Method == "POST" || r.Method == "PATCH" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "OPTIONS"
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

var ErrMalformedRequestLine = fmt.Errorf("malformed request-line")
var ErrUnsupportedHTTPVersion = fmt.Errorf("unsupported http version")
var SEPERATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, b, ErrMalformedRequestLine
	}
	startLine := b[:idx]
	restOfMessage := b[idx+len(SEPERATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMessage, ErrMalformedRequestLine
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMessage, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	if !rl.ValidMethod() {
		return nil, restOfMessage, ErrMalformedRequestLine
	}

	return rl, restOfMessage, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to io.ReadAll(): %w", err)
	}

	str := string(data)
	rl, _, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, err
}
