package tui

import (
	"testing"

	"github.com/Benjamin-Connelly/lookit/internal/index"
)

func TestTreeLess(t *testing.T) {
	tests := []struct {
		name string
		a, b index.FileEntry
		want bool
	}{
		{
			"dir before file at same level",
			index.FileEntry{RelPath: "docs", IsDir: true},
			index.FileEntry{RelPath: "README.md"},
			true,
		},
		{
			"file after dir at same level",
			index.FileEntry{RelPath: "README.md"},
			index.FileEntry{RelPath: "docs", IsDir: true},
			false,
		},
		{
			"alphabetical files",
			index.FileEntry{RelPath: "a.md"},
			index.FileEntry{RelPath: "b.md"},
			true,
		},
		{
			"alphabetical reverse",
			index.FileEntry{RelPath: "b.md"},
			index.FileEntry{RelPath: "a.md"},
			false,
		},
		{
			"case insensitive",
			index.FileEntry{RelPath: "Apple.md"},
			index.FileEntry{RelPath: "banana.md"},
			true,
		},
		{
			"parent before child",
			index.FileEntry{RelPath: "docs", IsDir: true},
			index.FileEntry{RelPath: "docs/guide.md"},
			true,
		},
		{
			"child after parent",
			index.FileEntry{RelPath: "docs/guide.md"},
			index.FileEntry{RelPath: "docs", IsDir: true},
			false,
		},
		{
			"nested dir before nested file",
			index.FileEntry{RelPath: "docs/api", IsDir: true},
			index.FileEntry{RelPath: "docs/readme.md"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := treeLess(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("treeLess(%q, %q) = %v, want %v", tt.a.RelPath, tt.b.RelPath, got, tt.want)
			}
		})
	}
}

func TestFileListModel_ListLen(t *testing.T) {
	m := FileListModel{
		visible:  make([]treeNode, 5),
		filtered: make([]index.FileEntry, 3),
	}

	// Not filtering — use visible
	if m.listLen() != 5 {
		t.Errorf("expected 5, got %d", m.listLen())
	}

	// Filtering — use filtered
	m.filtering = true
	if m.listLen() != 3 {
		t.Errorf("expected 3, got %d", m.listLen())
	}

	// Filter set but not in active filtering mode
	m.filtering = false
	m.filter = "query"
	if m.listLen() != 3 {
		t.Errorf("expected 3 (filter active), got %d", m.listLen())
	}
}

func TestFileListModel_VisibleRows(t *testing.T) {
	m := FileListModel{height: 30}
	if m.visibleRows() != 30 {
		t.Errorf("expected 30, got %d", m.visibleRows())
	}

	m.filtering = true
	if m.visibleRows() != 28 { // -2 for filter input
		t.Errorf("expected 28, got %d", m.visibleRows())
	}

	m.height = 0
	if m.visibleRows() < 1 {
		t.Error("should default to at least 1 row")
	}
}

func TestFileListModel_MoveUpDown(t *testing.T) {
	m := FileListModel{
		visible: make([]treeNode, 10),
		height:  20,
	}

	m.MoveDown()
	if m.cursor != 1 {
		t.Errorf("cursor should be 1, got %d", m.cursor)
	}

	// Move to end
	for i := 0; i < 20; i++ {
		m.MoveDown()
	}
	if m.cursor != 9 {
		t.Errorf("cursor should clamp at 9, got %d", m.cursor)
	}

	m.MoveUp()
	if m.cursor != 8 {
		t.Errorf("cursor should be 8, got %d", m.cursor)
	}

	// Move to top
	for i := 0; i < 20; i++ {
		m.MoveUp()
	}
	if m.cursor != 0 {
		t.Error("cursor should clamp at 0")
	}
}

func TestFileListModel_StartClearFilter(t *testing.T) {
	entries := []index.FileEntry{{RelPath: "a.md"}, {RelPath: "b.md"}}
	m := FileListModel{
		entries:  entries,
		filtered: entries,
	}

	m.StartFilter()
	if !m.filtering {
		t.Error("should be in filter mode")
	}
	if m.filter != "" {
		t.Error("filter should be empty")
	}

	m.filter = "something"
	m.cursor = 1
	m.ClearFilter()
	if m.filtering {
		t.Error("should exit filter mode")
	}
	if m.filter != "" {
		t.Error("filter should be cleared")
	}
	if m.cursor != 0 {
		t.Error("cursor should reset to 0")
	}
	if len(m.filtered) != 2 {
		t.Error("filtered should be reset to all entries")
	}
}

func TestFileListModel_Selected_Empty(t *testing.T) {
	m := FileListModel{}
	if m.Selected() != nil {
		t.Error("Selected on empty list should return nil")
	}
}

func TestFileListModel_Selected_ClampsCursor(t *testing.T) {
	entries := []index.FileEntry{{RelPath: "a.md"}}
	m := FileListModel{
		filtered: entries,
		cursor:   5, // out of bounds
	}
	sel := m.Selected()
	if sel == nil || sel.RelPath != "a.md" {
		t.Error("should clamp cursor and return entry")
	}
}

func TestFileListModel_SelectedVisible_Empty(t *testing.T) {
	m := FileListModel{}
	if m.SelectedVisible() != nil {
		t.Error("SelectedVisible on empty list should return nil")
	}
}

func TestFileListModel_SelectedVisible_UsesFiltered(t *testing.T) {
	entries := []index.FileEntry{{RelPath: "a.md"}}
	m := FileListModel{
		filtering: true,
		filtered:  entries,
	}
	sel := m.SelectedVisible()
	if sel == nil || sel.RelPath != "a.md" {
		t.Error("SelectedVisible in filter mode should use filtered list")
	}
}

func TestFileListModel_IsAncestorCollapsed(t *testing.T) {
	m := FileListModel{
		collapsed: map[string]bool{
			"docs": true,
		},
	}

	if !m.isAncestorCollapsed("docs/guide.md") {
		t.Error("docs/guide.md should have collapsed ancestor 'docs'")
	}
	if m.isAncestorCollapsed("src/main.go") {
		t.Error("src/main.go has no collapsed ancestor")
	}
	if m.isAncestorCollapsed("README.md") {
		t.Error("root file has no collapsed ancestor")
	}
}

func TestFileListModel_MoveDown_Scrolling(t *testing.T) {
	m := FileListModel{
		visible: make([]treeNode, 30),
		height:  10,
	}

	// Move past visible area
	for i := 0; i < 15; i++ {
		m.MoveDown()
	}

	// Offset should have scrolled
	if m.offset == 0 {
		t.Error("offset should have scrolled")
	}
	// Cursor should be visible
	visible := m.visibleRows()
	if m.cursor < m.offset || m.cursor >= m.offset+visible {
		t.Errorf("cursor %d should be in visible range [%d, %d)",
			m.cursor, m.offset, m.offset+visible)
	}
}

func TestFileListModel_MoveUp_Scrolling(t *testing.T) {
	m := FileListModel{
		visible: make([]treeNode, 30),
		height:  10,
		cursor:  15,
		offset:  10,
	}

	for i := 0; i < 10; i++ {
		m.MoveUp()
	}

	if m.cursor < m.offset {
		t.Error("cursor should be >= offset after scrolling up")
	}
}
