// cmd/hoho-tui/main.go
package main

import (
	"fmt"
	"os"

	"hoho/client"
	"hoho/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hoho-tui <port>")
		os.Exit(1)
	}
	port := os.Args[1]

	conn, err := client.Connect(port)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	p := tea.NewProgram(ui.New(conn), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
