package main

// This is a simple client for the go-chat-server.

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	go func() {
		reader := bufio.NewReader(conn)
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if err != nil {
				break
			}
			fmt.Print(string(buf[:n]))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		if strings.TrimSpace(message) == "/quit" {
			break
		}

		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			break
		}
	}
}
