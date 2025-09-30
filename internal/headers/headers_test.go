package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderParse(t *testing.T) {

	tests := []struct {
		name    string
		data    string
		headers []string
		values  []string
		length  int
		done    bool
		wantErr bool
	}{
		{
			name:    "valid single header",
			headers: []string{"Host"},
			data:    "Host: localhost:42069\r\n\r\n",
			length:  25,
			done:    true,
			wantErr: false,
			values:  []string{"localhost:42069"},
		}, {
			name:    "invalid spacing header",
			wantErr: true,
			data:    "     Host  : localhost:42069\r\n\r\n",
			length:  0,
			done:    false,
		}, {
			name:    "valid single header with extra white space",
			headers: []string{"Host"},
			length:  46,
			data:    "         Host: localhost:42069            \r\n\r\n",
			done:    true,
			wantErr: false,
			values:  []string{"localhost:42069"},
		}, {
			name:    "valide two headers",
			wantErr: false,
			data:    "Host: localhost:42069\r\n    Foo: bar  \r\n\r\n",
			length:  41,
			done:    true,
			headers: []string{"Host", "Foo"},
			values:  []string{"localhost:42069", "bar"},
		}, {
			name:    "malformed header",
			wantErr: true,
			data:    "H©st: localhost:42069\r\n\r\n",
			length:  0,
			done:    false,
		}, {
			name:    "valide multiple values",
			wantErr: false,
			data:    "Host: localhost:42069\r\nSet-Person: abhishek\r\nSet-Person: prithvi\r\nSet-Person: singh\r\n\r\n",
			length:  87,
			done:    true,
			headers: []string{"Host", "Set-Person"},
			values:  []string{"localhost:42069", "abhishek,prithvi,singh"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := NewHeaders()
			data := []byte(tt.data)
			n, done, err := headers.Parse(data)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error got none")
			}
			for i := range tt.headers {
				if value, ok := headers.Get(tt.headers[i]); ok {
					assert.Equal(t, tt.values[i], value)
				} else {
					t.Fatalf("expected value: %s found for none", tt.values[i])
				}
			}
			assert.Equal(t, tt.length, n)
			assert.Equal(t, tt.done, done)
		})
	}
}
