package styles

import "github.com/charmbracelet/lipgloss"

var Url = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#6c44b3"))

var Key = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#4078c0"))

var Header = lipgloss.NewStyle().
	PaddingLeft(2)

var Tag = lipgloss.NewStyle().
	Background(lipgloss.Color("#42474f")).
	Foreground(lipgloss.Color("#fff")).
	PaddingLeft(1).
	PaddingRight(1)
