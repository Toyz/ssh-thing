package components

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#0a7bca")).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 2)

	TabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#444444"))

	TabGap = lipgloss.NewStyle().
		Background(lipgloss.Color("#333333"))

	ViewportStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#0a7bca")).
			Padding(1, 2)

	ScrollableStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFB000")).
			Padding(1, 2)

	ScrolledUpStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF8800")).
			BorderBottom(true).
			BorderBottomForeground(lipgloss.Color("#FF5500")).
			Padding(1, 2)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	HelpShortKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87D7FF"))

	HelpShortDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A7B9C9"))

	HelpShortSeparator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#5F5F5F"))

	HelpFullKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87D7FF"))

	HelpFullDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A7B9C9"))

	HelpFullSeparator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#5F5F5F"))

	ScrollUpIndicator   = "↑"
	ScrollDownIndicator = "↓"
)
