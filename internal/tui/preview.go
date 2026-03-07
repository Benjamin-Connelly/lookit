package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PreviewModel renders file content in the preview pane.
type PreviewModel struct {
	content       string
	lines         []string
	filePath      string
	scroll        int
	width         int
	height        int
	highlightLine int // -1 = no highlight
}

// NewPreviewModel creates a preview pane.
func NewPreviewModel() PreviewModel {
	return PreviewModel{highlightLine: -1}
}

// SetContent updates the preview with rendered content.
func (m *PreviewModel) SetContent(path, content string) {
	m.filePath = path
	m.content = content
	m.lines = strings.Split(content, "\n")
	m.scroll = 0
	m.highlightLine = -1
}

// ScrollUp scrolls the preview up.
func (m *PreviewModel) ScrollUp(lines int) {
	m.scroll -= lines
	if m.scroll < 0 {
		m.scroll = 0
	}
}

// ScrollDown scrolls the preview down.
func (m *PreviewModel) ScrollDown(lines int) {
	m.scroll += lines
	maxScroll := m.maxScroll()
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}
}

// ScrollToBottom scrolls to the end of content.
func (m *PreviewModel) ScrollToBottom() {
	m.scroll = m.maxScroll()
}

func (m *PreviewModel) maxScroll() int {
	max := len(m.lines) - m.height
	if max < 0 {
		return 0
	}
	return max
}

// View renders the preview content.
func (m PreviewModel) View() string {
	if m.content == "" {
		placeholder := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)
		return placeholder.Render("Select a file to preview")
	}

	if len(m.lines) == 0 {
		return m.content
	}

	end := m.scroll + m.height
	if end > len(m.lines) {
		end = len(m.lines)
	}
	start := m.scroll
	if start >= len(m.lines) {
		start = len(m.lines) - 1
	}
	if start < 0 {
		start = 0
	}

	visible := make([]string, end-start)
	copy(visible, m.lines[start:end])

	// Highlight the line with the active link cursor
	if m.highlightLine >= start && m.highlightLine < end {
		idx := m.highlightLine - start
		hlStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("81"))
		visible[idx] = hlStyle.Render("▶ " + visible[idx])
	}

	result := strings.Join(visible, "\n")

	// Scroll position indicator if content exceeds viewport
	if len(m.lines) > m.height && m.height > 0 {
		pct := 0
		maxS := m.maxScroll()
		if maxS > 0 {
			pct = m.scroll * 100 / maxS
		}
		indicator := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		result += "\n" + indicator.Render(fmt.Sprintf("%s %d%%", strings.Repeat("─", 10), pct))
	}

	return result
}
