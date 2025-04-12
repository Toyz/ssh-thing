package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type ScrollView struct {
	viewport     viewport.Model
	buffer       *ScrollBuffer
	userScrolled bool
	wordWrap     bool
	width        int
	height       int
}

func NewScrollView() *ScrollView {
	vp := viewport.New(0, 0)
	vp.Style = ViewportStyle

	return &ScrollView{
		viewport:     vp,
		buffer:       NewScrollBuffer(DefaultMaxLines),
		userScrolled: false,
		wordWrap:     false,
	}
}

func (s *ScrollView) SetContent(content string) {
	s.viewport.SetContent(content)
}

func (s *ScrollView) Append(content string) {
	s.buffer.Append(content)
}

func (s *ScrollView) Clear() {
	s.buffer.Clear()
	s.viewport.SetContent("")
	s.userScrolled = false
}

func (s *ScrollView) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.viewport.Width = width
	s.viewport.Height = height
}

func (s *ScrollView) GotoBottom() {
	s.viewport.GotoBottom()
}

func (s *ScrollView) UserScrolled() bool {
	return s.userScrolled
}

func (s *ScrollView) SetUserScrolled(scrolled bool) {
	s.userScrolled = scrolled
}

func (s *ScrollView) ResetUserScrolledIfAtBottom() {
	if s.viewport.AtBottom() {
		s.userScrolled = false
	}
}

func (s *ScrollView) LineCount() int {
	return s.buffer.LineCount()
}

func (s *ScrollView) IsScrollable() bool {
	return s.viewport.Height < s.buffer.LineCount()
}

func (s *ScrollView) ToggleWordWrap() {
	s.wordWrap = !s.wordWrap
	content := s.buffer.String()
	if s.wordWrap {
		content = s.wrapContent(content)
	}
	s.viewport.SetContent(content)
}

func (s *ScrollView) IsWordWrapped() bool {
	return s.wordWrap
}

func (s *ScrollView) wrapContent(content string) string {
	if s.width <= 0 {
		return content
	}

	effectiveWidth := s.width - 4

	if effectiveWidth <= 10 {
		return content
	}

	lines := strings.Split(content, "\n")
	var wrappedLines []string

	for _, line := range lines {
		if len(line) <= effectiveWidth {
			wrappedLines = append(wrappedLines, line)
			continue
		}

		currentPos := 0
		for currentPos < len(line) {
			endPos := currentPos + effectiveWidth

			if endPos >= len(line) {
				wrappedLines = append(wrappedLines, line[currentPos:])
				break
			}

			lastSpace := strings.LastIndex(line[currentPos:endPos], " ")

			if lastSpace != -1 {
				wrappedLines = append(wrappedLines, line[currentPos:currentPos+lastSpace])
				currentPos += lastSpace + 1
			} else {
				wrappedLines = append(wrappedLines, line[currentPos:endPos])
				currentPos = endPos
			}
		}
	}

	return strings.Join(wrappedLines, "\n")
}

func (s *ScrollView) UpdateContent(colorizer func(string) string) {
	content := s.buffer.String()
	if colorizer != nil {
		content = colorizer(content)
	}
	if s.wordWrap {
		content = s.wrapContent(content)
	}
	s.viewport.SetContent(content)
}

func (s *ScrollView) ViewportModel() *viewport.Model {
	return &s.viewport
}

func (s *ScrollView) View() string {
	if s.IsScrollable() {
		if s.userScrolled {
			s.viewport.Style = ScrolledUpStyle
		} else {
			s.viewport.Style = ScrollableStyle
		}

		hasMoreAbove := s.viewport.YOffset > 0
		hasMoreBelow := !s.viewport.AtBottom()

		if hasMoreAbove || hasMoreBelow {
			var indicator string
			if hasMoreAbove && hasMoreBelow {
				indicator = " " + ScrollUpIndicator + ScrollDownIndicator + " "
			} else if hasMoreAbove {
				indicator = " " + ScrollUpIndicator + " "
			} else if hasMoreBelow {
				indicator = " " + ScrollDownIndicator + " "
			}

			styledIndicator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFB000")).
				Bold(true).
				Render(indicator)

			content := s.viewport.View()
			lines := strings.Split(content, "\n")

			if len(lines) > 0 {
				lastLineIdx := len(lines) - 1
				lastLine := lines[lastLineIdx]

				indicatorPlacement := lipgloss.PlaceHorizontal(
					s.viewport.Width,
					lipgloss.Right,
					styledIndicator,
				)

				lines[lastLineIdx] = lipgloss.PlaceHorizontal(
					s.viewport.Width,
					lipgloss.Left,
					lastLine,
				)

				return strings.Join(lines, "\n") + "\n" + indicatorPlacement
			}

			return content
		}
	} else {
		s.viewport.Style = ViewportStyle
	}

	return s.viewport.View()
}
