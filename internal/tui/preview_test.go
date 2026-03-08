package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func newTestPreview(lines int, height int) *PreviewModel {
	content := make([]string, lines)
	for i := range content {
		content[i] = strings.Repeat("x", 20)
	}
	m := NewPreviewModel()
	m.SetContent("test.md", strings.Join(content, "\n"))
	m.height = height
	m.scrolloff = 3
	return &m
}

func TestPreviewModel_NewEmpty(t *testing.T) {
	m := NewPreviewModel()
	if m.highlightLine != -1 {
		t.Errorf("highlight should be -1, got %d", m.highlightLine)
	}
	if m.cursorLine != 0 {
		t.Error("cursor should start at 0")
	}
}

func TestPreviewModel_SetContent(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("foo.md", "line1\nline2\nline3")
	if len(m.lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(m.lines))
	}
	if m.filePath != "foo.md" {
		t.Errorf("expected path foo.md, got %q", m.filePath)
	}
	if m.scroll != 0 {
		t.Error("scroll should reset to 0")
	}
	if m.visualMode {
		t.Error("visual mode should be cleared")
	}
}

func TestPreviewModel_SetContent_ResetsSearch(t *testing.T) {
	m := NewPreviewModel()
	m.searchMode = true
	m.searchQuery = "test"
	m.SetContent("bar.md", "new content")
	if m.searchMode {
		t.Error("search mode should be cleared on SetContent")
	}
	if m.searchQuery != "" {
		t.Error("search query should be cleared on SetContent")
	}
}

func TestPreviewModel_CursorDown(t *testing.T) {
	m := newTestPreview(50, 20)

	m.CursorDown()
	if m.cursorLine != 1 {
		t.Errorf("cursor should be 1, got %d", m.cursorLine)
	}

	// Move to bottom margin threshold
	for i := 0; i < 50; i++ {
		m.CursorDown()
	}
	if m.cursorLine != 49 {
		t.Errorf("cursor should clamp at 49, got %d", m.cursorLine)
	}
}

func TestPreviewModel_CursorUp(t *testing.T) {
	m := newTestPreview(50, 20)
	m.cursorLine = 5
	m.scroll = 2

	m.CursorUp()
	if m.cursorLine != 4 {
		t.Errorf("cursor should be 4, got %d", m.cursorLine)
	}

	// Move to top
	for i := 0; i < 10; i++ {
		m.CursorUp()
	}
	if m.cursorLine != 0 {
		t.Error("cursor should clamp at 0")
	}
}

func TestPreviewModel_CursorTo(t *testing.T) {
	m := newTestPreview(50, 20)

	m.CursorTo(25)
	if m.cursorLine != 25 {
		t.Errorf("expected cursor at 25, got %d", m.cursorLine)
	}
	// Should be visible
	if m.cursorLine < m.scroll || m.cursorLine >= m.scroll+m.height {
		t.Error("cursor should be within visible range after CursorTo")
	}

	// Clamp negative
	m.CursorTo(-5)
	if m.cursorLine != 0 {
		t.Errorf("negative line should clamp to 0, got %d", m.cursorLine)
	}

	// Clamp past end
	m.CursorTo(999)
	if m.cursorLine != 49 {
		t.Errorf("past-end should clamp to last line, got %d", m.cursorLine)
	}
}

func TestPreviewModel_ScrollUpDown(t *testing.T) {
	m := newTestPreview(100, 20)

	m.ScrollDown(10)
	if m.scroll != 10 {
		t.Errorf("scroll should be 10, got %d", m.scroll)
	}

	m.ScrollUp(5)
	if m.scroll != 5 {
		t.Errorf("scroll should be 5, got %d", m.scroll)
	}

	// Clamp at bottom
	m.ScrollDown(1000)
	max := len(m.lines) - m.height
	if m.scroll != max {
		t.Errorf("scroll should clamp at %d, got %d", max, m.scroll)
	}

	// Clamp at top
	m.ScrollUp(1000)
	if m.scroll != 0 {
		t.Error("scroll should clamp at 0")
	}
}

func TestPreviewModel_ScrollToBottom(t *testing.T) {
	m := newTestPreview(100, 20)
	m.ScrollToBottom()
	expected := len(m.lines) - m.height
	if m.scroll != expected {
		t.Errorf("scroll should be %d, got %d", expected, m.scroll)
	}
}

func TestPreviewModel_MaxScroll_ShortContent(t *testing.T) {
	m := newTestPreview(5, 20) // fewer lines than height
	m.ScrollDown(10)
	if m.scroll != 0 {
		t.Error("short content should have maxScroll=0")
	}
}

func TestPreviewModel_VisualMode(t *testing.T) {
	m := newTestPreview(50, 20)
	m.cursorLine = 10

	m.EnterVisualMode()
	if !m.visualMode {
		t.Error("should be in visual mode")
	}
	if m.visualAnchor != 10 || m.visualStart != 10 || m.visualEnd != 10 {
		t.Error("visual anchor/start/end should all be 10")
	}

	// Extend down
	m.VisualCursorDown()
	m.VisualCursorDown()
	if m.visualStart != 10 || m.visualEnd != 12 {
		t.Errorf("expected range 10-12, got %d-%d", m.visualStart, m.visualEnd)
	}

	// Extend up past anchor
	for i := 0; i < 5; i++ {
		m.VisualCursorUp()
	}
	if m.visualStart != 7 || m.visualEnd != 10 {
		t.Errorf("expected range 7-10, got %d-%d", m.visualStart, m.visualEnd)
	}

	m.ExitVisualMode()
	if m.visualMode {
		t.Error("should have exited visual mode")
	}
}

func TestPreviewModel_SelectedSourceLines_NoVisual(t *testing.T) {
	m := newTestPreview(50, 20)
	m.scroll = 5
	start, end := m.SelectedSourceLines()
	if start != 6 || end != 6 { // 1-based: scroll+1
		t.Errorf("expected 6,6 got %d,%d", start, end)
	}
}

func TestPreviewModel_SelectedSourceLines_Visual(t *testing.T) {
	m := newTestPreview(50, 20)
	m.cursorLine = 10
	m.EnterVisualMode()
	m.VisualCursorDown()
	m.VisualCursorDown()

	start, end := m.SelectedSourceLines()
	if start != 11 || end != 13 { // 1-based
		t.Errorf("expected 11,13 got %d,%d", start, end)
	}
}

func TestPreviewModel_ReadingGuide(t *testing.T) {
	m := NewPreviewModel()
	if m.readingGuide {
		t.Error("reading guide should default off")
	}
	m.ToggleReadingGuide()
	if !m.readingGuide {
		t.Error("reading guide should be on")
	}
	m.ToggleReadingGuide()
	if m.readingGuide {
		t.Error("reading guide should be off")
	}
}

func TestPreviewModel_SetSourceInfo(t *testing.T) {
	m := NewPreviewModel()
	m.SetSourceInfo(100, true)
	if m.sourceLineCount != 100 || !m.isCodeFile {
		t.Error("SetSourceInfo should store values")
	}
}

func TestPreviewModel_GutterWidth(t *testing.T) {
	tests := []struct {
		lines int
		want  int
	}{
		{5, 2},    // single digit: "N "
		{99, 3},   // two digits: "NN "
		{999, 4},  // three digits
		{9999, 5}, // four digits
	}
	for _, tt := range tests {
		m := newTestPreview(tt.lines, 20)
		got := m.gutterWidth()
		if got != tt.want {
			t.Errorf("gutterWidth(%d lines) = %d, want %d", tt.lines, got, tt.want)
		}
	}
}

// --- Search tests ---

func TestPreviewModel_SearchMode(t *testing.T) {
	m := newTestPreview(20, 10)

	m.EnterSearchMode()
	if !m.searchMode {
		t.Error("should be in search mode")
	}
	if m.searchQuery != "" {
		t.Error("search query should be empty")
	}
	if m.searchHistIdx != -1 {
		t.Error("search history index should be -1")
	}

	m.SearchInput('h')
	m.SearchInput('e')
	if m.searchQuery != "he" {
		t.Errorf("query should be 'he', got %q", m.searchQuery)
	}

	m.SearchBackspace()
	if m.searchQuery != "h" {
		t.Errorf("query should be 'h' after backspace, got %q", m.searchQuery)
	}

	m.SearchBackspace()
	if m.searchQuery != "" {
		t.Error("query should be empty after two backspaces")
	}

	// Backspace on empty is safe
	m.SearchBackspace()
	if m.searchQuery != "" {
		t.Error("extra backspace should be no-op")
	}
}

func TestPreviewModel_SearchMatches(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("test.md", "Hello world\nfoo bar\nhello again\nbaz")
	m.height = 20

	m.EnterSearchMode()
	m.SearchInput('h')
	m.SearchInput('e')
	m.SearchInput('l')
	m.SearchInput('l')
	m.SearchInput('o')

	if len(m.searchMatches) != 2 {
		t.Errorf("expected 2 matches for 'hello', got %d", len(m.searchMatches))
	}
	if m.searchMatches[0] != 0 || m.searchMatches[1] != 2 {
		t.Errorf("expected matches at lines 0,2, got %v", m.searchMatches)
	}
}

func TestPreviewModel_SearchNextPrev(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("test.md", "match\nno\nmatch\nno\nmatch")
	m.height = 20

	m.EnterSearchMode()
	for _, r := range "match" {
		m.SearchInput(r)
	}

	if len(m.searchMatches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(m.searchMatches))
	}

	// Start at 0
	if m.searchCurrent != 0 {
		t.Errorf("should start at match 0, got %d", m.searchCurrent)
	}

	m.NextMatch()
	if m.searchCurrent != 1 {
		t.Errorf("expected match 1, got %d", m.searchCurrent)
	}

	m.NextMatch()
	if m.searchCurrent != 2 {
		t.Errorf("expected match 2, got %d", m.searchCurrent)
	}

	// Wrap around
	m.NextMatch()
	if m.searchCurrent != 0 {
		t.Error("should wrap to match 0")
	}

	m.PrevMatch()
	if m.searchCurrent != 2 {
		t.Error("should wrap to last match")
	}
}

func TestPreviewModel_SearchNoMatches(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("test.md", "hello world")
	m.height = 20
	m.EnterSearchMode()
	for _, r := range "zzzzz" {
		m.SearchInput(r)
	}
	if len(m.searchMatches) != 0 {
		t.Error("should have no matches")
	}
	// NextMatch/PrevMatch should be safe with no matches
	m.NextMatch()
	m.PrevMatch()
}

func TestPreviewModel_SearchRegex(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("test.md", "foo123\nbar456\nbaz")
	m.height = 20

	m.EnterSearchMode()
	m.ToggleSearchRegex()
	if !m.searchRegex {
		t.Error("regex mode should be on")
	}

	for _, r := range `\d+` {
		m.SearchInput(r)
	}

	if len(m.searchMatches) != 2 {
		t.Errorf("expected 2 regex matches, got %d", len(m.searchMatches))
	}

	m.ToggleSearchRegex()
	// After toggle back to substring, `\d+` is literal — no matches
	if len(m.searchMatches) != 0 {
		t.Errorf("expected 0 literal matches for '\\d+', got %d", len(m.searchMatches))
	}
}

func TestPreviewModel_SearchRegex_Invalid(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("test.md", "hello")
	m.height = 20
	m.EnterSearchMode()
	m.ToggleSearchRegex()
	for _, r := range "[invalid" {
		m.SearchInput(r)
	}
	// Invalid regex should not panic, just return no matches
	if len(m.searchMatches) != 0 {
		t.Error("invalid regex should produce 0 matches")
	}
}

func TestPreviewModel_SearchHistory(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("test.md", "hello world foo bar")
	m.height = 20

	// First search
	m.EnterSearchMode()
	for _, r := range "hello" {
		m.SearchInput(r)
	}
	m.ExitSearchMode()

	// Second search
	m.EnterSearchMode()
	for _, r := range "world" {
		m.SearchInput(r)
	}
	m.ExitSearchMode()

	if len(m.searchHistory) != 2 {
		t.Fatalf("expected 2 history entries, got %d", len(m.searchHistory))
	}
	if m.searchHistory[0] != "world" {
		t.Error("most recent should be 'world'")
	}

	// Browse history
	m.EnterSearchMode()
	m.SearchHistoryUp()
	if m.searchQuery != "world" {
		t.Errorf("expected 'world', got %q", m.searchQuery)
	}
	m.SearchHistoryUp()
	if m.searchQuery != "hello" {
		t.Errorf("expected 'hello', got %q", m.searchQuery)
	}
	// Clamp at end
	m.SearchHistoryUp()
	if m.searchQuery != "hello" {
		t.Error("should clamp at oldest entry")
	}

	m.SearchHistoryDown()
	if m.searchQuery != "world" {
		t.Errorf("expected 'world', got %q", m.searchQuery)
	}
	m.SearchHistoryDown()
	if m.searchQuery != "" {
		t.Error("should return to empty query at bottom of history")
	}
}

func TestPreviewModel_SearchHistory_Dedup(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("test.md", "hello")
	m.height = 20

	for i := 0; i < 3; i++ {
		m.EnterSearchMode()
		for _, r := range "hello" {
			m.SearchInput(r)
		}
		m.ExitSearchMode()
	}

	if len(m.searchHistory) != 1 {
		t.Errorf("duplicate searches should be deduped, got %d entries", len(m.searchHistory))
	}
}

func TestPreviewModel_SearchHistoryEmpty(t *testing.T) {
	m := NewPreviewModel()
	m.SetContent("test.md", "hello")
	m.height = 20
	m.EnterSearchMode()
	// No history — up/down should be safe
	m.SearchHistoryUp()
	m.SearchHistoryDown()
	if m.searchQuery != "" {
		t.Error("query should remain empty with no history")
	}
}

func TestPreviewModel_IsSearchMatch(t *testing.T) {
	m := NewPreviewModel()
	m.searchMatches = []int{2, 5, 8}
	m.searchCurrent = 1

	if !m.isSearchMatch(5) {
		t.Error("line 5 should be a match")
	}
	if m.isSearchMatch(3) {
		t.Error("line 3 should not be a match")
	}
	if !m.isCurrentSearchMatch(5) {
		t.Error("line 5 should be the current match")
	}
	if m.isCurrentSearchMatch(2) {
		t.Error("line 2 is a match but not current")
	}
}

func TestPreviewModel_IsCurrentSearchMatch_OutOfBounds(t *testing.T) {
	m := NewPreviewModel()
	m.searchCurrent = -1
	if m.isCurrentSearchMatch(0) {
		t.Error("should return false for negative searchCurrent")
	}
	m.searchCurrent = 5
	m.searchMatches = []int{1, 2}
	if m.isCurrentSearchMatch(1) {
		t.Error("should return false when searchCurrent is out of bounds")
	}
}

func TestHighlightSearchInLine(t *testing.T) {
	// The function uses lipgloss styles, so we can't assert exact output.
	// Verify it doesn't lose content.
	result := highlightSearchInLine("Hello World Hello", "hello", lipgloss.NewStyle())
	if !strings.Contains(result, "Hello") && !strings.Contains(result, "hello") {
		t.Error("result should contain the matched text")
	}
	if !strings.Contains(result, "World") {
		t.Error("result should contain non-matched text")
	}
}

func TestHighlightSearchInLine_NoMatch(t *testing.T) {
	original := "no match here"
	result := highlightSearchInLine(original, "zzz", lipgloss.NewStyle())
	if result != original {
		t.Errorf("no-match should return original line, got %q", result)
	}
}

func TestPreviewModel_CursorScrolloff(t *testing.T) {
	m := newTestPreview(50, 20)
	m.scrolloff = 5

	// Move cursor down past the scrolloff margin
	for i := 0; i < 20; i++ {
		m.CursorDown()
	}
	// After 20 moves: cursor=20, scroll should keep cursor within scrolloff of bottom
	if m.cursorLine < m.scroll+m.scrolloff {
		t.Errorf("cursor %d should be at least scrolloff(%d) from top scroll(%d)",
			m.cursorLine, m.scrolloff, m.scroll)
	}
}
