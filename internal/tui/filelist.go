package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Benjamin-Connelly/lookit/internal/index"
)

// FileListModel manages the file list panel with fuzzy filtering.
type FileListModel struct {
	idx       *index.Index
	entries   []index.FileEntry
	filtered  []index.FileEntry
	cursor    int
	filter    string
	filtering bool
	offset    int // scroll offset
	height    int
}

// NewFileListModel creates a file list panel.
func NewFileListModel(idx *index.Index) FileListModel {
	entries := idx.Entries()
	return FileListModel{
		idx:      idx,
		entries:  entries,
		filtered: entries,
	}
}

// Selected returns the currently selected entry, or nil if empty.
func (m *FileListModel) Selected() *index.FileEntry {
	if len(m.filtered) == 0 {
		return nil
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	return &m.filtered[m.cursor]
}

// SetFilter updates the fuzzy filter query and refilters.
func (m *FileListModel) SetFilter(query string) {
	m.filter = query
	if query == "" {
		m.filtered = m.entries
	} else {
		m.filtered = m.idx.FuzzySearch(query, 100)
	}
	m.cursor = 0
	m.offset = 0
}

// StartFilter enters filter mode.
func (m *FileListModel) StartFilter() {
	m.filtering = true
	m.filter = ""
}

// ClearFilter exits filter mode and resets the list.
func (m *FileListModel) ClearFilter() {
	m.filtering = false
	m.filter = ""
	m.filtered = m.entries
	m.cursor = 0
	m.offset = 0
}

// MoveUp moves the cursor up.
func (m *FileListModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
	}
}

// MoveDown moves the cursor down.
func (m *FileListModel) MoveDown() {
	if m.cursor < len(m.filtered)-1 {
		m.cursor++
		visible := m.visibleRows()
		if visible > 0 && m.cursor >= m.offset+visible {
			m.offset = m.cursor - visible + 1
		}
	}
}

func (m *FileListModel) visibleRows() int {
	rows := m.height
	if m.filtering {
		rows -= 2 // filter input + separator
	}
	if rows < 1 {
		rows = 20
	}
	return rows
}

// View renders the file list.
func (m FileListModel) View() string {
	var b strings.Builder

	if m.filtering {
		filterStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
		b.WriteString(filterStyle.Render("/ "+m.filter+"_") + "\n")
		b.WriteString(strings.Repeat("─", 20) + "\n")
	}

	if len(m.filtered) == 0 {
		dim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		if m.filter != "" {
			b.WriteString(dim.Render("No matches for: " + m.filter))
		} else {
			b.WriteString(dim.Render("No files found"))
		}
		return b.String()
	}

	visible := m.visibleRows()
	end := m.offset + visible
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("62")).
		Bold(true)
	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))
	dirStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("75"))
	mdStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("114"))

	for i := m.offset; i < end; i++ {
		entry := m.filtered[i]
		icon := "  "
		if entry.IsDir {
			icon = "📁"
		} else if entry.IsMarkdown {
			icon = "📝"
		}

		line := fmt.Sprintf(" %s %s", icon, entry.RelPath)

		if i == m.cursor {
			b.WriteString(cursorStyle.Render("▸"+line) + "\n")
		} else if entry.IsDir {
			b.WriteString(dirStyle.Render(" "+line) + "\n")
		} else if entry.IsMarkdown {
			b.WriteString(mdStyle.Render(" "+line) + "\n")
		} else {
			b.WriteString(normalStyle.Render(" "+line) + "\n")
		}
	}

	// Scroll indicator
	if len(m.filtered) > visible {
		pct := 0
		if len(m.filtered)-visible > 0 {
			pct = m.offset * 100 / (len(m.filtered) - visible)
		}
		indicator := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		b.WriteString(indicator.Render(fmt.Sprintf(" %d/%d (%d%%)", m.cursor+1, len(m.filtered), pct)))
	}

	return b.String()
}
