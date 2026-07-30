package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	raddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Error resolving address")
		return
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		fmt.Println("error dialing UDP: ", err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unable to establish connection")
			return
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println("error writing to UDP connection: ", err)
			return
		}

	}
}
