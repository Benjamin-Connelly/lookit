package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCommandPalette_New(t *testing.T) {
	p := NewCommandPalette()
	if p.IsActive() {
		t.Error("new palette should not be active")
	}
}

func TestCommandPalette_OpenClose(t *testing.T) {
	p := NewCommandPalette()
	p.RegisterCommand(CommandEntry{Name: "test", Description: "desc"})

	p.Open()
	if !p.IsActive() {
		t.Error("should be active after Open")
	}
	if p.input != "" {
		t.Error("input should be cleared on Open")
	}
	if len(p.filtered) != 1 {
		t.Error("all commands should be in filtered on Open")
	}

	p.Close()
	if p.IsActive() {
		t.Error("should not be active after Close")
	}
}

func TestCommandPalette_SetInput_Filter(t *testing.T) {
	p := NewCommandPalette()
	p.RegisterCommand(CommandEntry{Name: "quit", Description: "exit"})
	p.RegisterCommand(CommandEntry{Name: "theme dark", Description: "dark theme"})
	p.RegisterCommand(CommandEntry{Name: "theme light", Description: "light theme"})
	p.Open()

	// Filter to "theme"
	p.SetInput("theme")
	if len(p.filtered) != 2 {
		t.Errorf("expected 2 matches for 'theme', got %d", len(p.filtered))
	}

	// Filter to "dark"
	p.SetInput("dark")
	if len(p.filtered) != 1 {
		t.Errorf("expected 1 match for 'dark', got %d", len(p.filtered))
	}

	// Empty restores all
	p.SetInput("")
	if len(p.filtered) != 3 {
		t.Errorf("expected 3 commands with empty filter, got %d", len(p.filtered))
	}
}

func TestCommandPalette_SetInput_CaseInsensitive(t *testing.T) {
	p := NewCommandPalette()
	p.RegisterCommand(CommandEntry{Name: "Quit", Description: "exit"})
	p.Open()

	p.SetInput("quit")
	if len(p.filtered) != 1 {
		t.Error("filter should be case-insensitive")
	}

	p.SetInput("QUIT")
	if len(p.filtered) != 1 {
		t.Error("uppercase filter should match lowercase name")
	}
}

func TestCommandPalette_MoveUpDown(t *testing.T) {
	p := NewCommandPalette()
	p.RegisterCommand(CommandEntry{Name: "a"})
	p.RegisterCommand(CommandEntry{Name: "b"})
	p.RegisterCommand(CommandEntry{Name: "c"})
	p.Open()

	if p.cursor != 0 {
		t.Error("cursor should start at 0")
	}

	p.MoveDown()
	p.MoveDown()
	if p.cursor != 2 {
		t.Errorf("expected cursor 2, got %d", p.cursor)
	}

	// Clamp at bottom
	p.MoveDown()
	if p.cursor != 2 {
		t.Error("should clamp at bottom")
	}

	p.MoveUp()
	if p.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", p.cursor)
	}

	// Clamp at top
	p.MoveUp()
	p.MoveUp()
	if p.cursor != 0 {
		t.Error("should clamp at top")
	}
}

func TestCommandPalette_Execute(t *testing.T) {
	called := false
	p := NewCommandPalette()
	p.RegisterCommand(CommandEntry{
		Name: "test",
		Action: func() tea.Msg {
			called = true
			return StatusMsg{Text: "done"}
		},
	})
	p.Open()

	msg := p.Execute()
	if !called {
		t.Error("action should be called")
	}
	if p.IsActive() {
		t.Error("palette should close after execute")
	}
	if sm, ok := msg.(StatusMsg); !ok || sm.Text != "done" {
		t.Errorf("unexpected message: %v", msg)
	}
}

func TestCommandPalette_Execute_NoAction(t *testing.T) {
	p := NewCommandPalette()
	p.RegisterCommand(CommandEntry{Name: "noop"})
	p.Open()

	msg := p.Execute()
	if msg != nil {
		t.Error("nil action should return nil message")
	}
	if p.IsActive() {
		t.Error("palette should close after execute")
	}
}

func TestCommandPalette_Execute_Empty(t *testing.T) {
	p := NewCommandPalette()
	p.Open()
	// No commands registered
	msg := p.Execute()
	if msg != nil {
		t.Error("should return nil for empty palette")
	}
}

func TestCommandPalette_SetInput_ResetsCursor(t *testing.T) {
	p := NewCommandPalette()
	p.RegisterCommand(CommandEntry{Name: "a"})
	p.RegisterCommand(CommandEntry{Name: "b"})
	p.Open()

	p.MoveDown()
	if p.cursor != 1 {
		t.Error("cursor should be 1")
	}

	p.SetInput("a")
	if p.cursor != 0 {
		t.Error("SetInput should reset cursor to 0")
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s, substr string
		want      bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "WORLD", true},
		{"Hello World", "lo wo", true},
		{"Hello", "Hello World", false}, // substr longer
		{"abc", "xyz", false},
		{"", "", true},
		{"abc", "", true},
		{"", "a", false},
	}
	for _, tt := range tests {
		got := containsIgnoreCase(tt.s, tt.substr)
		if got != tt.want {
			t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
		}
	}
}

func TestCommandPalette_View_Active(t *testing.T) {
	p := NewCommandPalette()
	p.RegisterCommand(CommandEntry{Name: "quit", Description: "exit"})
	p.Open()
	p.input = "qu"

	view := p.View()
	if view == "" {
		t.Error("view should not be empty when active")
	}
}

func TestCommandPalette_View_Inactive(t *testing.T) {
	p := NewCommandPalette()
	if p.View() != "" {
		t.Error("view should be empty when inactive")
	}
}
