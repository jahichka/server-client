package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
	"uuid"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

// ANSI escape codes for text colors
var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
}

var clients = make(map[string]net.Conn) // Map to store connections using UUID
var mutex = &sync.Mutex{}

func handleClient(conn net.Conn) {
	clientUUID := uuid.New().String()
	rand.Seed(time.Now().UnixNano())
	chosenColor := colors[rand.Intn(len(colors))]
	conn.Write([]byte(chosenColor + "\n"))

	fmt.Println("New client connected:", clientUUID)

	mutex.Lock()
	clients[clientUUID] = conn
	mutex.Unlock()

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
			response = fmt.Sprintf("Server time: %s \n", time.Now().Format(time.RFC1123))
		case strings.HasPrefix(command, "ECHO"):
			response = "Echo: " + strings.TrimSpace(command[5:]) + "\n"
		case strings.HasPrefix(command, "CHAT"):
			defer func() {
				mutex.Lock()
				delete(clients, clientUUID)
				mutex.Unlock()
				conn.Close()
			}()
			reader := bufio.NewReader(conn)
			for {
				message, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Client", clientUUID, "left.")
					return
				}

				// Broadcast the message to all clients except the sender
				for id, c := range clients {
					if id != clientUUID {
						c.Write([]byte("[" + clientUUID + "]: " + message))
					}
				}
			}
		case strings.HasPrefix(command, "EXIT"):
			response = "Goodbye!\n"
			conn.Write([]byte(response))
			return
		default:
			response = "Unknown command\n"
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
