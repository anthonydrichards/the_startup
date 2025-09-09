package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders_Parse(t *testing.T) {
	// Test: Valid single header
	headers := make(Headers)
	data := []byte("HOST: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	headerValue, _ := headers.Get("host")
	assert.Equal(t, "localhost:42069", headerValue)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = make(Headers)
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header without colon
	headers = make(Headers)
	data = []byte("Host localhost\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid multiple headers
	headers = make(Headers)
	data = []byte("Host: localhost:42069\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n")
	totalParsed := 0
	for {
		n, done, err = headers.Parse(data[totalParsed:])
		require.NoError(t, err)
		totalParsed += n
		if done {
			break
		}
	}
	assert.Equal(t, 3, len(headers))
	headerValue, _ = headers.Get("host")
	assert.Equal(t, "localhost:42069", headerValue)
	headerValue, _ = headers.Get("user-agent")
	assert.Equal(t, "Go-http-client/1.1", headerValue)
	headerValue, _ = headers.Get("accept-encoding")
	assert.Equal(t, "gzip", headerValue)
	headerValue, ok := headers.Get("non-existent-header")
	assert.False(t, ok)
	assert.Equal(t, "", headerValue)
	assert.Equal(t, 80, totalParsed)

	// Test: Incomplete header (no CRLF)
	headers = make(Headers)
	data = []byte("Host: localhost:42069")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header name
	headers = make(Headers)
	data = []byte("H@st: localhost\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple headers with same name
	headers = make(Headers)
	data = []byte("Set-Cookie: id=123\r\nSet-Cookie: token=abc\r\n\r\n")
	totalParsed = 0
	for {
		n, done, err = headers.Parse(data[totalParsed:])
		require.NoError(t, err)
		totalParsed += n
		if done {
			break
		}
	}
	assert.Equal(t, 1, len(headers))
	headerValue, _ = headers.Get("set-cookie")
	assert.Equal(t, "id=123, token=abc", headerValue)
	assert.Equal(t, 45, totalParsed)

}
