package main

import (
	"errors"
	"fmt"
	"github.com/danielwolber-wood/httpfromtcp/internal/request"
	"io"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		buffer := []byte{}
		for {
			chunk := make([]byte, 8)
			i, err := f.Read(chunk)
			if err != nil && !errors.Is(err, io.EOF) {
				fmt.Printf("Error while reading from %v: %v\n", f, err)
				return
			}
			for index := 0; index < i; index++ {
				char := chunk[index]
				if string(char) == "\n" {
					ch <- string(buffer)
					buffer = []byte{}
				} else {
					buffer = append(buffer, char)
				}
			}
			if i == 0 && errors.Is(err, io.EOF) {
				// i.e., after the entire file is read, send whatever is left in the buffer and close the channel
				if len(buffer) != 0 {
					ch <- string(buffer)
				}
				close(ch)
				return
			}
		}
	}()
	return ch
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		fmt.Printf("Error opening listener: %v\n", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v", err)
			return
		}
		fmt.Println("Connection accepted")
		/*readChan := getLinesChannel(conn)
		for item := range readChan {
			fmt.Println(item)
		} */
		data, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error reading header")
		}
		fmt.Printf("Data: %v\n", data)
		fmt.Printf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n", data.RequestLine.Method, data.RequestLine.RequestTarget, data.RequestLine.HttpVersion)
		fmt.Println("Connection closed")
	}
}
