package request

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestLineParse(t *testing.T) {

	// Test: Good GET Request line
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)

	// Test: Good POST Request line
	r, err = RequestFromReader(strings.NewReader("POST /coffee/hello HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee/hello", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Missing HTTP Version
	_, err = RequestFromReader(strings.NewReader("POST /coffee/hello \r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)

	// Test: Empty Request
	_, err = RequestFromReader(strings.NewReader(""))
	require.Error(t, err)

	// Test: Out of order request line
	_, err = RequestFromReader(strings.NewReader("/coffee/hello POST HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)

	// Test: Out of order request line
	_, err = RequestFromReader(strings.NewReader("/coffee/hello HTTP/1.1 GET\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)

	// Test: Invalid Method
	_, err = RequestFromReader(strings.NewReader("SMASH /coffee/hello HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)

	// Test: Good GET Request line with chunkReader
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path with chunkReader
	reader = &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

}

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
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func TestRequestLineParseWithChunkReader(t *testing.T) {
	testCases := []struct {
		name            string
		request         string
		expectError     bool
		expectedMethod  string
		expectedTarget  string
		expectedVersion string
	}{
		{
			name:            "Good GET Request line",
			request:         "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError:     false,
			expectedMethod:  "GET",
			expectedTarget:  "/",
			expectedVersion: "1.1",
		},
		{
			name:            "Good GET Request line with path",
			request:         "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError:     false,
			expectedMethod:  "GET",
			expectedTarget:  "/coffee",
			expectedVersion: "1.1",
		},
		{
			name:        "Invalid number of parts in request line",
			request:     "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
		{
			name:            "Good POST Request line",
			request:         "POST /coffee/hello HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError:     false,
			expectedMethod:  "POST",
			expectedTarget:  "/coffee/hello",
			expectedVersion: "1.1",
		},
		{
			name:        "Missing HTTP Version",
			request:     "POST /coffee/hello \r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
		{
			name:        "Empty Request",
			request:     "",
			expectError: true,
		},
		{
			name:        "Out of order request line (target first)",
			request:     "/coffee/hello POST HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
		{
			name:        "Out of order request line (version second)",
			request:     "/coffee/hello HTTP/1.1 GET\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
		{
			name:        "Invalid Method",
			request:     "SMASH /coffee/hello HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
	}

	// Test with different chunk sizes for each test case
	for _, tc := range testCases {
		// Calculate chunk sizes to test
		requestLen := len(tc.request)
		chunkSizes := []int{
			1,              // Single byte at a time
			3,              // Small chunks
			5,              // Medium chunks
			requestLen / 2, // Half the request
			requestLen,     // Entire request
		}

		// Add some interesting chunk sizes that might reveal edge cases
		if requestLen > 10 {
			chunkSizes = append(chunkSizes, 10)
		}
		if requestLen > 20 {
			chunkSizes = append(chunkSizes, 20)
		}

		// Run each test case with different chunk sizes
		for _, chunkSize := range chunkSizes {
			testName := fmt.Sprintf("%s with chunk size %d", tc.name, chunkSize)

			t.Run(testName, func(t *testing.T) {
				reader := &chunkReader{
					data:            tc.request,
					numBytesPerRead: chunkSize,
				}

				r, err := RequestFromReader(reader)

				if tc.expectError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.NotNil(t, r)
					assert.Equal(t, tc.expectedMethod, r.RequestLine.Method)
					assert.Equal(t, tc.expectedTarget, r.RequestLine.RequestTarget)
					assert.Equal(t, tc.expectedVersion, r.RequestLine.HttpVersion)
				}
			})
		}
	}
}

// Alternative implementation that combines the original test with chunked reading tests
func TestRequestLineParseComprehensive(t *testing.T) {
	testCases := []struct {
		name            string
		request         string
		expectError     bool
		expectedMethod  string
		expectedTarget  string
		expectedVersion string
	}{
		{
			name:            "Good GET Request line",
			request:         "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError:     false,
			expectedMethod:  "GET",
			expectedTarget:  "/",
			expectedVersion: "1.1",
		},
		{
			name:            "Good GET Request line with path",
			request:         "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError:     false,
			expectedMethod:  "GET",
			expectedTarget:  "/coffee",
			expectedVersion: "1.1",
		},
		{
			name:        "Invalid number of parts in request line",
			request:     "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
		{
			name:            "Good POST Request line",
			request:         "POST /coffee/hello HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError:     false,
			expectedMethod:  "POST",
			expectedTarget:  "/coffee/hello",
			expectedVersion: "1.1",
		},
		{
			name:        "Missing HTTP Version",
			request:     "POST /coffee/hello \r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
		{
			name:        "Empty Request",
			request:     "",
			expectError: true,
		},
		{
			name:        "Out of order request line (target first)",
			request:     "/coffee/hello POST HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
		{
			name:        "Out of order request line (version second)",
			request:     "/coffee/hello HTTP/1.1 GET\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
		{
			name:        "Invalid Method",
			request:     "SMASH /coffee/hello HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectError: true,
		},
	}

	// First run tests with regular string reader (as in original test)
	for _, tc := range testCases {
		t.Run(tc.name+" with StringReader", func(t *testing.T) {
			r, err := RequestFromReader(strings.NewReader(tc.request))

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, r)
				assert.Equal(t, tc.expectedMethod, r.RequestLine.Method)
				assert.Equal(t, tc.expectedTarget, r.RequestLine.RequestTarget)
				assert.Equal(t, tc.expectedVersion, r.RequestLine.HttpVersion)
			}
		})
	}

	// Then run tests with chunk reader using various sizes
	for _, tc := range testCases {
		requestLen := len(tc.request)
		if requestLen == 0 {
			// Skip empty request for multiple chunk sizes
			t.Run(tc.name+" with ChunkReader (size=1)", func(t *testing.T) {
				reader := &chunkReader{
					data:            tc.request,
					numBytesPerRead: 1,
				}
				_, err := RequestFromReader(reader)
				require.Error(t, err)
			})
			continue
		}

		// Test different chunk sizes
		chunkSizes := []int{
			1,              // Single byte
			3,              // Small chunks
			requestLen / 4, // Quarter of request
			requestLen / 2, // Half of request
			requestLen,     // Entire request
		}

		// Add boundary test for CRLF edge cases
		if requestLen > 4 {
			chunkSizes = append(chunkSizes, 4) // For CRLF boundary tests
		}

		for _, chunkSize := range chunkSizes {
			testName := fmt.Sprintf("%s with ChunkReader (size=%d)", tc.name, chunkSize)

			t.Run(testName, func(t *testing.T) {
				reader := &chunkReader{
					data:            tc.request,
					numBytesPerRead: chunkSize,
				}

				r, err := RequestFromReader(reader)

				if tc.expectError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.NotNil(t, r)
					assert.Equal(t, tc.expectedMethod, r.RequestLine.Method)
					assert.Equal(t, tc.expectedTarget, r.RequestLine.RequestTarget)
					assert.Equal(t, tc.expectedVersion, r.RequestLine.HttpVersion)
				}
			})
		}
	}
}
