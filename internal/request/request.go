package request

import (
	"bytes"
	"fmt"
	"httpServerGo/internal/headers"
	"io"
)

type parserState string

const (
	StateInit   parserState = "init"
	StateHeader parserState = "header"
	StateBody   parserState = "body"
	StateDone   parserState = "done"
	StateErr    parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) ValidMethod() bool {
	return r.Method == "GET" || r.Method == "POST" || r.Method == "PATCH" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "OPTIONS"
}

type RequestBody struct{
  Body string
  BodyLength string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	state       parserState
	Body        RequestBody
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		currentData := data[read:]

		switch r.state {
		case StateErr:
			return 0, ErroRequestInErrorState
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
      if err == ErrIncomplete {
        break outer
      }
			if err != nil {
				r.state = StateErr
				return 0, err
			}
			r.RequestLine = *rl
			read += n

			r.state = StateHeader
		case StateHeader:

			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateErr
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				r.state = StateDone
			}
		case StateBody:
      
		case StateDone:
			break outer
		default:
			panic("i have doemsome shitty mistake")
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
var ErrIncomplete = errors.New("incomplete data")
var SEPERATOR = []byte("\r\n")

func parseRequestBody(b []byte,contentLength int) (string,error) {
  if contentLength == -1 {
    return "",fmt.Errorf("content length not provided")
  }
  if 
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, ErrIncomplete
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
