---
name: tui-design
description: Use when designing, building, or modifying a Text User Interface (TUI) or command-line interface in Go.
---

# TUI Design (Charm)

## Overview
When tasked with building or designing a Text User Interface (TUI) in Go, you must use the [Charm](https://charm.sh/) ecosystem of libraries. The goal is to produce interfaces that are colorful, vivid, modern, and highly interactive, breaking away from traditional monochrome terminal output.

## When to Use
- Whenever the user asks to build a TUI, CLI application, or terminal tool.
- When tasked with designing interactive command-line components.
- When writing Go applications that run in the terminal and require user interaction.

## Core Pattern
Build upon the Elm architecture using Charm's specialized libraries:
1.  **Bubble Tea (`github.com/charmbracelet/bubbletea`):** The core framework. Use this for the Model, Update (event handling), and View (rendering) loop.
2.  **Lip Gloss (`github.com/charmbracelet/lipgloss`):** The styling engine. Use this extensively for vivid, colorful, and modern styling. Define styles for borders, margins, padding, foregrounds, and backgrounds. Never use raw ANSI escape codes.
3.  **Bubbles (`github.com/charmbracelet/bubbles`):** Ready-to-use UI components. Leverage existing components like text inputs, spinners, progress bars, tables, and lists instead of building them from scratch.
4.  **Huh (`github.com/charmbracelet/huh`):** (If applicable) Use for building modern, interactive forms and prompts quickly.

## Design Philosophy: Colorful, Vivid, Modern
-   **Color Palettes:** Choose a modern, vibrant color palette. Use Lip Gloss to apply these colors consistently across the UI. Consider dark mode/light mode adaptive colors if Lip Gloss's adaptive colors are appropriate.
-   **Layout:** Use Lip Gloss for layout management. Employ Flexbox-like concepts (using `JoinHorizontal`, `JoinVertical`, and alignments) to create clean, structured interfaces.
-   **Borders & Spacing:** Generously use borders (`lipgloss.RoundedBorder()`, `lipgloss.NormalBorder()`) and padding to separate logical areas of the UI. This contributes significantly to a "modern" feel.
-   **Interactivity:** Ensure the application responds fluidly to user input using Bubble Tea's `Update` function.

## Implementation Example
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define vivid, modern styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 2)
)

type model struct{}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	title := titleStyle.Render("Vivid TUI Design")
	content := boxStyle.Render("Press 'q' to quit this modern interface.")
	
	return lipgloss.JoinVertical(lipgloss.Left, title, content)
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
```

## Common Mistakes
-   Writing raw ANSI codes instead of using Lip Gloss.
-   Making the UI monochrome or ignoring styling entirely.
-   Not separating the Model, Update, and View logic properly according to Bubble Tea patterns.
-   Reinventing the wheel instead of using existing `charmbracelet/bubbles` components.