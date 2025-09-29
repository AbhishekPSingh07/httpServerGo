package request

import (
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
    if cr.pos >= len(cr.data) {
        return 0, io.EOF
    }
    
    // Read up to numBytesPerRead or remaining data, whichever is smaller
    remainingData := len(cr.data) - cr.pos
    bytesToRead := min(cr.numBytesPerRead, remainingData, len(p))
    
    n = copy(p, cr.data[cr.pos:cr.pos+bytesToRead])
    cr.pos += n
    
    return n, nil
}

func TestUnit_ParseRequestHeader(t *testing.T) {

	tests := []struct {
		name    string
		request string
		headers []string
		values  []string
		wantErr bool
	}{
		{
			name:    "Standard headers",
			request: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			headers: []string{"Host", "User-Agent", "Accept"},
			values:  []string{"localhost:42069", "curl/7.81.0", "*/*"},
			wantErr: false,
		},
		{
			name:    "Malformed Headers",
			request: "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			randomNumber := rng.Intn(7) + 2

			reader := &chunkReader{
				data:            tt.request,
				numBytesPerRead: randomNumber,
			}
			r, err := RequestFromReader(reader)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, but got one: %v", err)
			}
			for i := range tt.headers {
				assert.Equal(t, tt.values[i], r.Headers.Get(tt.headers[i]))
			}

		})
	}
}

func TestUnit_RequestFromReader(t *testing.T) {

	tests := []struct {
		name        string
		request     string
		wantMethod  string
		wantTarget  string
		wantVersion string
		wantErr     bool
		errExpected error
	}{
		{
			name:        "Good Get Request line",
			request:     "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			wantMethod:  "GET",
			wantTarget:  "/",
			wantVersion: "1.1",
			wantErr:     false,
		}, {
			name:        "Good Get Request line with path",
			request:     "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			wantMethod:  "GET",
			wantTarget:  "/coffee",
			wantVersion: "1.1",
		}, {
			name:    "Invalid Number of parts in request line",
			request: "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			wantErr: true,
		},
		{
			name:        "Good Post Request With Path",
			request:     "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Type: application/json\r\nContent-Length: 22\r\n\r\n{\"flavor\":\"dark mode\"}",
			wantMethod:  "POST",
			wantTarget:  "/coffee",
			wantVersion: "1.1",
		},
		{
			name:    "Invalid method out of order request line",
			request: "/coffee POST HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Type: application/json\r\nContent-Length: 22\r\n\r\n{\"flavor\":\"dark mode\"}",
			wantErr: true,
		},
		{
			name:    "Invalid method out of order request line",
			request: "POST /coffee HTTP/1.2\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Type: application/json\r\nContent-Length: 22\r\n\r\n{\"flavor\":\"dark mode\"}",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			randomNumber := rng.Intn(7) + 2

			reader := &chunkReader{
				data:            tt.request,
				numBytesPerRead: randomNumber,
			}
			r, err := RequestFromReader(reader)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, but got one: %v", err)
			}
			assert.Equal(t, tt.wantMethod, r.RequestLine.Method)
			assert.Equal(t, tt.wantTarget, r.RequestLine.RequestTarget)
			assert.Equal(t, tt.wantVersion, r.RequestLine.HttpVersion)
		})
	}
}
