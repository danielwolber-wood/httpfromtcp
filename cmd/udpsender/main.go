package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error on resolving: %v\n", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("Error on dialing: %v\n", err)
		return
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		s, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from stdin: %v\n", err)
			return
		}
		// I know this is wrong, but isn't a string just a slice of bytes?
		_, err = conn.Write([]byte(s))
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			return
		}
	}
}
