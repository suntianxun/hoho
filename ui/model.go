// ui/model.go
package ui

import (
	"fmt"
	"net"
	"strings"

	"hoho/client"

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
	Conn   net.Conn
	State  []client.VarState
}

func New(conn net.Conn) Model {
	return Model{
		Conn: conn,
	}
}

type stateMsg client.StateMessage
type errMsg error

func fetchState(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		if conn == nil {
			return errMsg(fmt.Errorf("no connection"))
		}
		msg, err := client.ReadState(conn)
		if err != nil {
			return errMsg(err)
		}
		return stateMsg(msg)
	}
}

func (m Model) Init() tea.Cmd {
	if m.Conn != nil {
		return fetchState(m.Conn)
	}
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
	case stateMsg:
		m.State = msg.Vars
		return m, fetchState(m.Conn)
	case errMsg:
		// Log or handle error if needed
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

	var varLines []string
	varLines = append(varLines, "Variables")
	varLines = append(varLines, "---------")
	for _, v := range m.State {
		// optionally skip private or format them
		line := fmt.Sprintf("%s (%s): %s", v.Name, v.Type, v.Value)
		varLines = append(varLines, line)
	}

	varsPane := paneStyle.Width(vw).Height(vh).Render(strings.Join(varLines, "\n"))
	termPane := paneStyle.Width(tw).Height(vh).Render("Terminal")
	
	return lipgloss.JoinHorizontal(lipgloss.Top, varsPane, termPane)
}
