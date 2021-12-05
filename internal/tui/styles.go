package tui

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#059396")).
			Padding(0, 1)

	selectedItemTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#00f9ff"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#00f9ff"}).
				Padding(0, 0, 0, 1)
	selectedItemDesc = selectedItemTitleStyle.Copy().
				Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#398082"})

	itemTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#059396"}).
			Padding(0, 0, 0, 2)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)
