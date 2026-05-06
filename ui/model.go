// ui/model.go
package ui

import (
	"fmt"
	"strings"

	"hoho/client"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	varsWidthPercentage = 40
	paddingAndBorders   = 4
)

var (
	paneStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).BorderForeground(lipgloss.Color("#874BFD"))
)

type Model struct {
	width  int
	height int
	Ready  bool
	Client *client.Client
	State  []client.VarState
	Err    error
}

func New(client *client.Client) Model {
	return Model{
		Client: client,
	}
}

type stateMsg client.StateMessage
type errMsg error

func fetchState(c *client.Client) tea.Cmd {
	return func() tea.Msg {
		if c == nil {
			return errMsg(fmt.Errorf("no connection"))
		}
		msg, err := c.ReadState()
		if err != nil {
			return errMsg(err)
		}
		return stateMsg(msg)
	}
}

func (m Model) Init() tea.Cmd {
	if m.Client != nil {
		return fetchState(m.Client)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "n" || msg.String() == "next" || msg.String() == "c" {
			if m.Client != nil {
				m.Client.SendCommand(msg.String())
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Ready = true
	case stateMsg:
		m.State = msg.Vars
		m.Err = nil
		return m, nil
	case errMsg:
		m.Err = error(msg)
	}
	return m, nil
}

func (m Model) View() string {
	if !m.Ready {
		return "Initializing..."
	}

	varsWidth := (m.width * varsWidthPercentage) / 100
	termWidth := m.width - varsWidth
	
	// Subtract borders and padding
	vw := varsWidth - paddingAndBorders
	if vw < 0 { vw = 0 }
	tw := termWidth - paddingAndBorders
	if tw < 0 { tw = 0 }
	vh := m.height - paddingAndBorders
	if vh < 0 { vh = 0 }

	var varLines []string
	if m.Err != nil {
		varLines = append(varLines, "Error: "+m.Err.Error())
	} else {
		varLines = append(varLines, "Variables")
		varLines = append(varLines, "---------")
		for _, v := range m.State {
			line := fmt.Sprintf("%s (%s): %s", v.Name, v.Type, v.Value)
			varLines = append(varLines, line)
		}
	}

	varsPane := paneStyle.Width(vw).Height(vh).Render(strings.Join(varLines, "\n"))
	termPane := paneStyle.Width(tw).Height(vh).Render("Terminal")
	
	return lipgloss.JoinHorizontal(lipgloss.Top, varsPane, termPane)
}
