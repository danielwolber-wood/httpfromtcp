package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	filepath := "./messages.txt"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error opening %v: %v\n", filepath, err)
	}
	for {
		arr := make([]byte, 8)
		i, err := file.Read(arr)
		if i == 0 && errors.Is(err, io.EOF) {
			os.Exit(0)
		} else if err != nil {
			fmt.Printf("Error while reading from %v: %v\n", filepath, err)
		}
		fmt.Printf("read: %s\n", arr)

	}

}
