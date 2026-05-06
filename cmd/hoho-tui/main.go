// cmd/hoho-tui/main.go
package main

import (
	"fmt"
	"os"
	"strconv"

	"hoho/client"
	"hoho/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hoho-tui <port>")
		os.Exit(1)
	}
	portStr := os.Args[1]

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		fmt.Println("Error: valid port number between 1 and 65535 is required")
		os.Exit(1)
	}

	c, err := client.Connect(portStr)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer c.Conn.Close()

	p := tea.NewProgram(ui.New(c), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
