package tui

import (
	"testing"
)

func TestStatusBarModel_New(t *testing.T) {
	m := NewStatusBarModel()
	if m.mode != "NORMAL" {
		t.Errorf("default mode should be NORMAL, got %q", m.mode)
	}
}

func TestStatusBarModel_SetFile(t *testing.T) {
	m := NewStatusBarModel()
	m.SetFile("README.md")
	if m.filePath != "README.md" {
		t.Errorf("expected path README.md, got %q", m.filePath)
	}
}

func TestStatusBarModel_SetMessage(t *testing.T) {
	m := NewStatusBarModel()
	m.SetMessage("Copied!")
	if m.message != "Copied!" {
		t.Errorf("expected message 'Copied!', got %q", m.message)
	}
}

func TestStatusBarModel_SetMode(t *testing.T) {
	m := NewStatusBarModel()
	m.SetMode("VISUAL")
	if m.mode != "VISUAL" {
		t.Errorf("expected mode VISUAL, got %q", m.mode)
	}
}

func TestStatusBarModel_ContextHints_Search(t *testing.T) {
	m := NewStatusBarModel()
	m.searchMode = true
	m.searchQuery = "hello"
	m.searchMatchCount = 3

	hints := m.contextHints()
	if hints == "" {
		t.Fatal("hints should not be empty")
	}
	if !containsStr(hints, "hello") {
		t.Error("should contain search query")
	}
	if !containsStr(hints, "3 matches") {
		t.Error("should contain match count")
	}
	if !containsStr(hints, "ctrl-r") {
		t.Error("should mention regex toggle")
	}
}

func TestStatusBarModel_ContextHints_Visual(t *testing.T) {
	m := NewStatusBarModel()
	m.visualMode = true
	m.visualRange = "L5-L10"

	hints := m.contextHints()
	if !containsStr(hints, "permalink") {
		t.Error("visual hints should mention permalink")
	}
	if !containsStr(hints, "L5-L10") {
		t.Error("should contain visual range")
	}
}

func TestStatusBarModel_ContextHints_Help(t *testing.T) {
	m := NewStatusBarModel()
	m.showingHelp = true
	hints := m.contextHints()
	if !containsStr(hints, "close help") {
		t.Error("help hints should mention close")
	}
}

func TestStatusBarModel_ContextHints_LinkActive(t *testing.T) {
	m := NewStatusBarModel()
	m.linkActive = true
	m.linkText = "[README.md]"
	hints := m.contextHints()
	if !containsStr(hints, "follow") {
		t.Error("link hints should mention follow")
	}
	if !containsStr(hints, "README.md") {
		t.Error("should contain link text")
	}
}

func TestStatusBarModel_ContextHints_ByFocus(t *testing.T) {
	tests := []struct {
		focus Panel
		want  string
	}{
		{PanelPreview, "search"},
		{PanelSide, "select"},
		{PanelFileList, "filter"},
	}
	for _, tt := range tests {
		m := NewStatusBarModel()
		m.focus = tt.focus
		hints := m.contextHints()
		if !containsStr(hints, tt.want) {
			t.Errorf("focus %d hints should contain %q, got %q", tt.focus, tt.want, hints)
		}
	}
}
