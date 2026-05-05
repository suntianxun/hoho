// ui/model.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	paneStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).BorderForeground(lipgloss.Color("#874BFD"))
)

type Model struct {
	Ready bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	varsPane := paneStyle.Width(40).Height(20).Render("Variables")
	termPane := paneStyle.Width(60).Height(20).Render("Terminal")
	
	return lipgloss.JoinHorizontal(lipgloss.Top, varsPane, termPane)
}
