package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaders_Parse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("    Host: localhost:42069  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// Test: 2 valid headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nBond: James\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[23:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "James", headers["Bond"])
	assert.Equal(t, 13, n)
	assert.False(t, done)

	// Test: 2 valid headers and parsing end of block
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nBond: James\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[23:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "James", headers["Bond"])
	assert.Equal(t, 13, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[36:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl/7.81.0", headers["User-Agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n a bunch of other stuff")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Empty(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

}
