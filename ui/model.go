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
	width  int
	height int
	Ready  bool
}

func InitialModel() Model {
	return Model{}
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Ready = true
	}
	return m, nil
}

func (m Model) View() string {
	if !m.Ready {
		return "Initializing..."
	}

	varsWidth := (m.width * 4) / 10
	termWidth := m.width - varsWidth
	
	// Subtract borders and padding (2 for border + 2 for padding = 4)
	vw := varsWidth - 4
	if vw < 0 { vw = 0 }
	tw := termWidth - 4
	if tw < 0 { tw = 0 }
	vh := m.height - 4
	if vh < 0 { vh = 0 }

	varsPane := paneStyle.Width(vw).Height(vh).Render("Variables")
	termPane := paneStyle.Width(tw).Height(vh).Render("Terminal")
	
	return lipgloss.JoinHorizontal(lipgloss.Top, varsPane, termPane)
}
