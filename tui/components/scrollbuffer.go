package components

import (
	"strings"
)

const (
	DefaultMaxLines = 1000
)

type ScrollBuffer struct {
	lines    []string
	maxLines int
}

func NewScrollBuffer(maxLines int) *ScrollBuffer {
	if maxLines <= 0 {
		maxLines = DefaultMaxLines
	}
	return &ScrollBuffer{
		lines:    make([]string, 0, maxLines),
		maxLines: maxLines,
	}
}

func (s *ScrollBuffer) Append(content string) {
	newLines := strings.Split(content, "\n")

	if len(s.lines) > 0 && !strings.HasSuffix(s.lines[len(s.lines)-1], "\n") && len(newLines) > 0 {
		s.lines[len(s.lines)-1] += newLines[0]
		newLines = newLines[1:]
	}

	s.lines = append(s.lines, newLines...)

	if len(s.lines) > s.maxLines {
		s.lines = s.lines[len(s.lines)-s.maxLines:]
	}
}

func (s *ScrollBuffer) Clear() {
	s.lines = make([]string, 0, s.maxLines)
}

func (s *ScrollBuffer) String() string {
	return strings.Join(s.lines, "\n")
}

func (s *ScrollBuffer) LineCount() int {
	return len(s.lines)
}
