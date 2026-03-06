package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestHighlight(t *testing.T) {
	r := NewCodeRenderer("dark", false)
	out, err := r.Highlight("main.go", `package main

func main() {
	println("hello")
}`)
	if err != nil {
		t.Fatalf("Highlight: %v", err)
	}
	// HTML output should contain class-based spans
	if !strings.Contains(out, "chroma") {
		t.Error("expected chroma CSS classes in HTML output")
	}
}

func TestHighlightTerminal(t *testing.T) {
	r := NewCodeRenderer("dark", true)
	out, err := r.Highlight("main.py", `print("hello")`)
	if err != nil {
		t.Fatalf("Highlight: %v", err)
	}
	// Terminal output should contain ANSI escapes
	if !strings.Contains(out, "\x1b[") {
		t.Error("expected ANSI escape codes in terminal output")
	}
}

func TestHighlightFallback(t *testing.T) {
	r := NewCodeRenderer("dark", false)
	// Unknown extension should still produce output
	out, err := r.Highlight("file.xyz123", "some content")
	if err != nil {
		t.Fatalf("Highlight: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output for unknown file type")
	}
}

func TestGetLanguage(t *testing.T) {
	r := NewCodeRenderer("dark", false)

	tests := []struct {
		file string
		want string
	}{
		{"main.go", "Go"},
		{"app.py", "Python"},
		{"index.js", "JavaScript"},
		{"unknown.xyz123", ""},
	}
	for _, tt := range tests {
		got := r.GetLanguage(tt.file)
		if got != tt.want {
			t.Errorf("GetLanguage(%q) = %q, want %q", tt.file, got, tt.want)
		}
	}
}

func TestCSS(t *testing.T) {
	r := NewCodeRenderer("dark", false)
	css, err := r.CSS()
	if err != nil {
		t.Fatalf("CSS: %v", err)
	}
	if !strings.Contains(css, ".chroma") {
		t.Error("CSS should contain .chroma selector")
	}
}

func TestSetTheme(t *testing.T) {
	r := NewCodeRenderer("dark", false)
	r.SetTheme("light")
	// Should use github style now
	css, _ := r.CSS()
	if css == "" {
		t.Error("CSS should not be empty after theme change")
	}
}

func TestHighlightToWriter(t *testing.T) {
	r := NewCodeRenderer("dark", false)
	var buf bytes.Buffer
	err := r.HighlightToWriter(&buf, "main.go", `package main`)
	if err != nil {
		t.Fatalf("HighlightToWriter: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestHighlightLines(t *testing.T) {
	r := NewCodeRenderer("dark", false)
	out, err := r.HighlightLines("main.go", "package main\n\nfunc main() {}\n", 2, 3)
	if err != nil {
		t.Fatalf("HighlightLines: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestListLanguages(t *testing.T) {
	r := NewCodeRenderer("dark", false)
	langs := r.ListLanguages()
	if len(langs) < 50 {
		t.Errorf("expected 50+ languages, got %d", len(langs))
	}
}
