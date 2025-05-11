package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	ParseState  int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParseState == 1 {
		return 0, fmt.Errorf("request is already parsed")
	}
	requestLine, bytesConsumed, err := parseRequestLine(string(data))
	if bytesConsumed == 0 && err != nil {
		return 0, err
	}
	if bytesConsumed != 0 {
		// request line was parsed correctly, so buffer should be shortened by bytesConsume
		// and then the number of bytes consumed from this request should be returned
		r.ParseState = 1
		r.RequestLine = requestLine
		return bytesConsumed, nil
	}
	return 0, nil
}

func parseRequestLine(line string) (RequestLine, int, error) {
	fmt.Printf("line: %v\n", line)

	if !strings.Contains(line, "\r\n") {
		return RequestLine{}, 0, nil
	}

	// We know the string contains a CRLF now, so we can find its index; that's how many bytes we consume
	index := strings.Index(line, "\r\n")
	bytesConsumed := len(line) - index + 2
	requestLine := line[:index]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid request line")
	}
	fmt.Printf("parts: %v\n", parts)

	validMethods := []string{"GET", "PUT", "POST", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT"}
	method := parts[0]
	if !slices.Contains(validMethods, method) {
		return RequestLine{}, 0, fmt.Errorf("invalid HTTP Method")
	}
	fmt.Printf("method: %v\n", method)

	if len(parts[2]) < 5 || len(parts[2]) > 8 {
		return RequestLine{}, 0, fmt.Errorf("invalid HTTP Version")
	}
	if !strings.Contains(parts[2], "/") {
		return RequestLine{}, 0, fmt.Errorf("invalid HTTP Version")
	}

	httpVersion := strings.Split(parts[2], "/")[1]
	validHTTPVersions := []string{"1.1", "2", "3"}
	if !slices.Contains(validHTTPVersions, httpVersion) {
		return RequestLine{}, 0, fmt.Errorf("invalid HTTP Version")
	}
	fmt.Printf("httpVersion: %v\n", httpVersion)
	requestTarget := parts[1]

	extractedRequestLine := RequestLine{HttpVersion: httpVersion, RequestTarget: requestTarget, Method: method}
	return extractedRequestLine, bytesConsumed, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 8)
	var allData []byte
	request := Request{ParseState: 0}
	for request.ParseState != 1 {
		// I don't know if I actually need bytesRead for anything
		bytesRead, err := reader.Read(buffer)
		if err != nil && !errors.Is(io.EOF, err) {
			return nil, err
		} else if errors.Is(io.EOF, err) && bytesRead == 0 && request.ParseState == 0 {
			// i.e., if we're at EOF, no bytes were read, and the parsing hasn't succeeded yet
			// the second conditional is needed because the first time we hit EOF, we may have read a partial buffer
			return nil, fmt.Errorf("EOF reached without successful parsing")
		}
		// doesn't this copy the entire allData buffer every time? That doesn't seem right
		allData = append(allData, buffer[:bytesRead]...)
		bytesConsumed, err := request.parse(allData)
		if err != nil {
			return nil, err
		}
		if bytesConsumed == 0 {
			// i.e., if no parsing happened
			continue
		}
		allData = allData[bytesConsumed:]
	}
	return &request, nil
}
