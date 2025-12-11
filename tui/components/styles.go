package components

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#bd93f9")). // Dracula Purple
			Foreground(lipgloss.Color("#282a36")). // Dracula Background
			Padding(0, 1)

	TabStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("#f8f8f2")). // Dracula Foreground
			Background(lipgloss.Color("#44475a"))  // Dracula Current Line

	TabGap = lipgloss.NewStyle().Background(lipgloss.Color("#282a36"))

	ViewportStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#bd93f9")) // Dracula Purple

	ScrollableStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ffb86c")) // Dracula Orange

	ScrolledUpStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff5555")) // Dracula Red

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282a36")). // Dark FG
			Background(lipgloss.Color("#bd93f9")). // Purple BG
			Padding(0, 1).
			Bold(true)

	StatusText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")).
			Background(lipgloss.Color("#44475a")).
			Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")). // Dracula Comment
			Background(lipgloss.Color("#44475a"))

	HelpStyleWithBorder = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6272a4")).
				Background(lipgloss.Color("#44475a"))

	HelpShortKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8be9fd")). // Dracula Cyan
			Background(lipgloss.Color("#44475a"))

	HelpShortDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")).
			Background(lipgloss.Color("#44475a"))

	HelpShortSeparator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#44475a")).
				Background(lipgloss.Color("#44475a"))

	HelpFullKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8be9fd")).
			Background(lipgloss.Color("#44475a"))

	HelpFullDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")).
			Background(lipgloss.Color("#44475a"))

	HelpFullSeparator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#44475a")).
				Background(lipgloss.Color("#44475a"))

	ScrollIndicatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffb86c")).
				Bold(true)

	ScrollUpIndicator   = "↑"
	ScrollDownIndicator = "↓"
)
