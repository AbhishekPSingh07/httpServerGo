package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")

func parseHeader(feildLine []byte) (string, string, error) {
	parts := bytes.SplitN(feildLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}
	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {

	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		// EMPTY HEADER
		if idx == 0 {
			done = true
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(rn)
		h[name] = value
	}

	return read, done, nil
}

func NewHeaders() Headers {
	return map[string]string{}
}
