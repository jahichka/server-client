package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

// ANSI escape codes for text colors and cursor movement
const (
	Blue      = "\033[34m"
	Reset     = "\033[0m"
	CursorUp  = "\033[1A"
	ClearLine = "\033[2K"
)

func listenForMessages(conn net.Conn, myColor string) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Disconnected from server.")
			conn.Close()
			os.Exit(1)
		}

		// Move the cursor up to overwrite the "Text to send:" prompt
		fmt.Print(CursorUp + ClearLine)

		// Use the color assigned by the server
		fmt.Println(myColor + message + Reset)

		// Reprint the prompt for the user
		fmt.Print("Text to send: ")
	}
}

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	myColor, _ := bufio.NewReader(conn).ReadString('\n')
	myColor = myColor[:len(myColor)-1]

	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command (TIME, ECHO <message>, CHAT, EXIT): ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		_, err := conn.Write([]byte(command))
		if err != nil {
			fmt.Println("Error sending command:", err)
			return
		}

		if strings.HasPrefix(command, "EXIT") {
			break
		}

		if strings.HasPrefix(command, "CHAT") {
			go listenForMessages(conn, myColor)

			inputReader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Text to send: ")
				input, err := inputReader.ReadString('\n')
				if err != nil {
					log.Println("Error reading input:", err)
					return
				}

				// Send the message to the server
				conn.Write([]byte(input))

				// Print the user's own message without color (on a new line)
				fmt.Print(CursorUp + ClearLine)
				fmt.Println("You: " + input[:len(input)-1]) // Remove newline from input
			}
		}
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			break
		}
		fmt.Println("Server response:", response)

	}
}
