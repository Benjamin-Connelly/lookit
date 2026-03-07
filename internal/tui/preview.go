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

	// Source line tracking (for permalink generation)
	sourceLineCount int  // total lines in source file
	isCodeFile      bool // true = rendered lines map 1:1 to source

	// Visual line selection
	visualMode   bool
	visualAnchor int // where selection started (fixed)
	visualStart  int // min(anchor, cursor)
	visualEnd    int // max(anchor, cursor)
	cursorLine   int // current cursor position in visual mode
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
	m.visualMode = false
	m.cursorLine = 0
}

// SetSourceInfo stores metadata about the source file for line mapping.
func (m *PreviewModel) SetSourceInfo(lineCount int, isCode bool) {
	m.sourceLineCount = lineCount
	m.isCodeFile = isCode
}

// EnterVisualMode starts line selection at the current scroll position.
func (m *PreviewModel) EnterVisualMode() {
	m.visualMode = true
	m.cursorLine = m.scroll
	m.visualAnchor = m.cursorLine
	m.visualStart = m.cursorLine
	m.visualEnd = m.cursorLine
}

// ExitVisualMode clears selection.
func (m *PreviewModel) ExitVisualMode() {
	m.visualMode = false
}

// VisualCursorDown moves the visual cursor down.
func (m *PreviewModel) VisualCursorDown() {
	if m.cursorLine < len(m.lines)-1 {
		m.cursorLine++
	}
	m.updateVisualRange()
	// Auto-scroll if cursor moves below viewport
	if m.cursorLine >= m.scroll+m.height {
		m.scroll = m.cursorLine - m.height + 1
	}
}

// VisualCursorUp moves the visual cursor up.
func (m *PreviewModel) VisualCursorUp() {
	if m.cursorLine > 0 {
		m.cursorLine--
	}
	m.updateVisualRange()
	// Auto-scroll if cursor moves above viewport
	if m.cursorLine < m.scroll {
		m.scroll = m.cursorLine
	}
}

func (m *PreviewModel) updateVisualRange() {
	if m.cursorLine < m.visualAnchor {
		m.visualStart = m.cursorLine
		m.visualEnd = m.visualAnchor
	} else {
		m.visualStart = m.visualAnchor
		m.visualEnd = m.cursorLine
	}
}

// SelectedSourceLines returns the 1-based source line range for the visual selection.
// For code files (1:1 mapping), this is exact. For markdown, it's approximate.
func (m *PreviewModel) SelectedSourceLines() (startLine, endLine int) {
	if !m.visualMode {
		// Single line at cursor/scroll position
		line := m.scroll + 1
		return line, line
	}
	return m.visualStart + 1, m.visualEnd + 1
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

	// Highlight visual selection range
	if m.visualMode {
		selStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("24")).
			Foreground(lipgloss.Color("255"))
		for i := range visible {
			lineIdx := start + i
			if lineIdx >= m.visualStart && lineIdx <= m.visualEnd {
				visible[i] = selStyle.Render(visible[i])
			}
		}
	}

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
