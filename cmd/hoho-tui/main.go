// cmd/hoho-tui/main.go
package main

import (
	"fmt"
	"os"

	"hoho/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
