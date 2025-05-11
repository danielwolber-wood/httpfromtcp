package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	fmt.Printf("data: %v\n", data)
	line := string(data)
	fmt.Printf("string(data): %v\n", string(data))
	if !strings.Contains(line, "\r\n") {
		return 0, false, nil
	}

	if strings.HasPrefix(line, "\r\n") {
		return 0, true, nil
	}

	index := strings.Index(line, "\r\n")
	fmt.Printf("index: %v\n", index)
	bytesConsumed := index + 2
	fmt.Printf("bytesConsumed: %v\n", bytesConsumed)
	headerLine := strings.TrimSpace(line[:index])

	if !strings.Contains(headerLine, ":") {
		return 0, false, fmt.Errorf("invalid header line: no colon")
	}
	splitIndex := strings.Index(headerLine, ":")
	key := strings.TrimLeft(headerLine[:splitIndex], " ")
	value := strings.TrimSpace(headerLine[splitIndex+1:])
	fmt.Printf("Key: %v\n", key)
	fmt.Printf("Value: %v\n", value)
	if strings.HasSuffix(key, " ") {
		return 0, false, fmt.Errorf("invalid header line: space preceding colon")
	}
	h[key] = value

	return bytesConsumed, false, nil

}
