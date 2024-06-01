package main

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
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command (TIME, ECHO <message>, EXIT): ")
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

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			break
		}

		fmt.Println("Server response:", response)

	}
}
