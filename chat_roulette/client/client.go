package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Client struct{}

func New() Client {
	return Client{}
}

func (c Client) Run() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	// Read from partner
	go func() {
		for {
			buffer := make([]byte, 1024)
			count, err := conn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s (%d)\n", buffer, count)
		}
	}()

	// Send to partner
	for scanner.Scan() {
		msg := scanner.Bytes()
		conn.Write(msg)
	}
}
