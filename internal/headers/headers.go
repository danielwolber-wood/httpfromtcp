package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func validateKey(key string) (bool, error) {
	allow := make(map[rune]struct{})
	chars := []rune{
		'!', '#', '$', '%', '&', '\'', '*',
		'+', '-', '.', '^', '_', '`', '|', '~', ':',
	}
	for i := '0'; i <= '9'; i++ {
		chars = append(chars, i)
	}
	for i := 'a'; i <= 'z'; i++ {
		chars = append(chars, i)
	}

	for _, char := range chars {
		allow[char] = struct{}{}
	}

	for _, r := range key {
		if _, ok := allow[r]; !ok {
			return false, nil
		}
	}
	return true, nil

}

func parseFieldLine(line string) (n int, key, value string, done bool, err error) {
	fmt.Printf("line: %v\n", line)
	if !strings.Contains(line, "\r\n") {
		return 0, "", "", false, nil
	}

	if strings.HasPrefix(line, "\r\n") {
		return 2, "", "", true, nil
	}

	index := strings.Index(line, "\r\n")
	fmt.Printf("index: %v\n", index)
	bytesConsumed := index + 2
	fmt.Printf("bytesConsumed: %v\n", bytesConsumed)
	headerLine := strings.TrimSpace(line[:index])

	if !strings.Contains(headerLine, ":") {
		return 0, "", "", false, fmt.Errorf("invalid header line: no colon")
	}
	var splitIndex int
	for i, v := range headerLine {
		if string(v) == ":" {
			if i > 0 && string(headerLine[i-1]) == " " {
				return 0, "", "", false, fmt.Errorf("invalid header line: space before colon")
			} else if i == 0 {
				return 0, "", "", false, fmt.Errorf("invalid header line: field begins with colon")
			}
			splitIndex = i
			break
		}
	}
	key = strings.ToLower(headerLine[:splitIndex])
	value = headerLine[splitIndex+1:] // +1 because I don't want it to contain the colon
	key = strings.TrimSpace(string(key))
	value = strings.TrimSpace(string(value))
	keyValid, _ := validateKey(key)
	if !keyValid {
		return 0, "", "", false, fmt.Errorf("invalid header line: invalid char in key")
	}
	return bytesConsumed, key, value, false, nil

}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Goal: this function is mostly concerned with the logistics of where one field begins and the next ends, and whether we're at the end of the field.
	// the actual parsing login is in parseFieldLine
	fmt.Printf("data: %v\n", data)
	fmt.Printf("string(data): %v\n", string(data))
	if !strings.Contains(string(data), "\r\n") {
		return 0, false, nil
	}
	// or should this be
	/// if string(data) == "\r\n\r\n"
	/// or strings.HasPrefix(string(data), "\r\n\r\n")
	if strings.HasPrefix(string(data), "\r\n") {
		return 2, true, nil
	}
	index := strings.Index(string(data), "\r\n")
	bytesConsumed, key, value, done, err := parseFieldLine(string(data)[:index+2])
	if err != nil {
		return 0, false, err
	}
	if bytesConsumed == 0 {
		return 0, false, nil
	}
	fmt.Printf("Key: %v\n ", key)
	fmt.Printf("Value: %v\n ", value)
	existingValue, ok := h[key]
	if ok {
		fmt.Printf("existingValue: %v\n", existingValue)
		result := fmt.Sprintf("%v, %v", existingValue, value)
		fmt.Printf("result: %v\n", result)
		h[key] = result

	} else {
		h[key] = value
	}
	return bytesConsumed, false, nil
}
