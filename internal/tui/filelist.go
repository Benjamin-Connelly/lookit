package tui

import (
	"github.com/Benjamin-Connelly/lookit/internal/index"
)

// FileListModel manages the file list panel with fuzzy filtering.
type FileListModel struct {
	idx      *index.Index
	entries  []index.FileEntry
	filtered []index.FileEntry
	cursor   int
	filter   string
	offset   int // scroll offset
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
	return &m.filtered[m.cursor]
}

// SetFilter updates the fuzzy filter query and refilters.
func (m *FileListModel) SetFilter(query string) {
	m.filter = query
	if query == "" {
		m.filtered = m.entries
	} else {
		m.filtered = m.idx.FuzzySearch(query)
	}
	m.cursor = 0
	m.offset = 0
}

// MoveUp moves the cursor up.
func (m *FileListModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown moves the cursor down.
func (m *FileListModel) MoveDown() {
	if m.cursor < len(m.filtered)-1 {
		m.cursor++
	}
}

// View renders the file list.
func (m FileListModel) View() string {
	if len(m.filtered) == 0 {
		return "No files found"
	}

	var s string
	for i, entry := range m.filtered {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		icon := " "
		if entry.IsDir {
			icon = " "
		} else if entry.IsMarkdown {
			icon = " "
		}

		s += cursor + icon + entry.RelPath + "\n"
	}
	return s
}
