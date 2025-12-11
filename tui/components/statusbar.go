package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	Width int
}

func NewStatusBar() *StatusBar {
	return &StatusBar{}
}

func (s *StatusBar) View(serverName string, status string, scrollPos string, helpView string) string {
	w := lipgloss.Width

	statusKey := StatusBarStyle.Render("STATUS")
	statusVal := StatusText.Copy().
		Foreground(statusColor(status)).
		Render(status)

	serverKey := StatusBarStyle.Render("SERVER")
	serverVal := StatusText.Render(serverName)

	scrollKey := StatusBarStyle.Render("SCROLL")
	scrollVal := StatusText.Render(scrollPos)

	// Add a small gap between sections using the background color of the text
	gap := StatusText.Render(" ")

	statusBlock := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey, statusVal, gap,
		serverKey, serverVal, gap,
		scrollKey, scrollVal, gap,
	)

	// Fill the rest of the width with the status bar background color
	availWidth := s.Width - w(statusBlock)
	if availWidth < 0 {
		availWidth = 0
	}
	
	helpVal := StatusText.Copy().
		Width(availWidth).
		Align(lipgloss.Right).
		Render(helpView)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		statusBlock,
		helpVal,
	)
}

func statusColor(status string) lipgloss.Color {
	switch strings.ToLower(status) {
	case "connected":
		return lipgloss.Color("#00FF00") // Green
	case "error":
		return lipgloss.Color("#FF0000") // Red
	case "connecting":
		return lipgloss.Color("#FFFF00") // Yellow
	default:
		return lipgloss.Color("#CCCCCC") // Grey
	}
}
