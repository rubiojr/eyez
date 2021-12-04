package main

import "github.com/charmbracelet/lipgloss"

var urlStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#6c44b3"))

var keyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#4078c0"))

var headersStyle = lipgloss.NewStyle().
	PaddingLeft(2)

var tagStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#42474f")).
	Foreground(lipgloss.Color("#fff")).
	PaddingLeft(1).
	PaddingRight(1)
