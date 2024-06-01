package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		command := string(buffer[:n])
		command = strings.TrimSpace(command)

		var response string

		switch {
		case strings.HasPrefix(command, "TIME"):
			response = fmt.Sprintf("Server time: %s", time.Now().Format(time.RFC1123))
		case strings.HasPrefix(command, "ECHO"):
			response = "Echo: " + strings.TrimSpace(command[5:])
		case strings.HasPrefix(command, "EXIT"):
			response = "Goodbye!"
			conn.Write([]byte(response))
			return
		default:
			response = "Unknown command"
		}

		conn.Write([]byte(response))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server is listening on port 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}
