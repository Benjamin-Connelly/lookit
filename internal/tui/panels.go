package tui

// PanelType identifies the type of side panel displayed.
type PanelType int

const (
	PanelNone PanelType = iota
	PanelTOC
	PanelBacklinks
	PanelGitInfo
	PanelBookmarks
)

// TOCEntry represents a heading in the table of contents.
type TOCEntry struct {
	Level int
	Text  string
	Line  int
}

// Bookmark represents a saved location.
type Bookmark struct {
	Path   string
	Title  string
	Scroll int
}

// SidePanelModel manages the optional side panel (TOC, backlinks, etc.).
type SidePanelModel struct {
	panelType PanelType
	toc       []TOCEntry
	bookmarks []Bookmark
	cursor    int
}

// NewSidePanelModel creates a side panel model.
func NewSidePanelModel() SidePanelModel {
	return SidePanelModel{
		panelType: PanelNone,
	}
}

// Toggle switches between showing and hiding the given panel type.
func (m *SidePanelModel) Toggle(pt PanelType) {
	if m.panelType == pt {
		m.panelType = PanelNone
	} else {
		m.panelType = pt
		m.cursor = 0
	}
}

// SetTOC updates the table of contents entries.
func (m *SidePanelModel) SetTOC(entries []TOCEntry) {
	m.toc = entries
}

// AddBookmark adds a bookmark.
func (m *SidePanelModel) AddBookmark(bm Bookmark) {
	m.bookmarks = append(m.bookmarks, bm)
}

// RemoveBookmark removes a bookmark by index.
func (m *SidePanelModel) RemoveBookmark(idx int) {
	if idx >= 0 && idx < len(m.bookmarks) {
		m.bookmarks = append(m.bookmarks[:idx], m.bookmarks[idx+1:]...)
	}
}

// Visible returns whether the side panel is currently shown.
func (m *SidePanelModel) Visible() bool {
	return m.panelType != PanelNone
}

// View renders the side panel.
func (m SidePanelModel) View() string {
	switch m.panelType {
	case PanelTOC:
		return m.viewTOC()
	case PanelBookmarks:
		return m.viewBookmarks()
	default:
		return ""
	}
}

func (m SidePanelModel) viewTOC() string {
	if len(m.toc) == 0 {
		return "No headings found"
	}
	var s string
	for i, entry := range m.toc {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		indent := ""
		for j := 1; j < entry.Level; j++ {
			indent += "  "
		}
		s += cursor + indent + entry.Text + "\n"
	}
	return s
}

func (m SidePanelModel) viewBookmarks() string {
	if len(m.bookmarks) == 0 {
		return "No bookmarks"
	}
	var s string
	for i, bm := range m.bookmarks {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		s += cursor + bm.Title + " (" + bm.Path + ")\n"
	}
	return s
}
