package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
	StateErr  parserState = "error"
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

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		switch r.state {
		case StateErr:
			return 0, ErroRequestInErrorState
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateErr
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateDone

		case StateDone:
			break outer
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateErr
}

var ErrMalformedRequestLine = fmt.Errorf("malformed request-line")
var ErrUnsupportedHTTPVersion = fmt.Errorf("unsupported http version")
var ErroRequestInErrorState = fmt.Errorf("request in error state")
var SEPERATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil 
	}
	startLine := b[:idx]
	read := idx + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	if !rl.ValidMethod() {
		return nil, 0, ErrMalformedRequestLine
	}

	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		// TODO : better handling of errs
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
