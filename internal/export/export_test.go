package export

import (
	"testing"
)

func TestTitleFromPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"README.md", "README"},
		{"my-doc.md", "my doc"},
		{"docs/guide.md", "guide"},
		{"notes.markdown", "notes"},
	}
	for _, tt := range tests {
		got := titleFromPath(tt.input)
		if got != tt.want {
			t.Errorf("titleFromPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestReplaceExt(t *testing.T) {
	tests := []struct {
		name, newExt, want string
	}{
		{"file.md", ".html", "file.html"},
		{"doc.markdown", ".html", "doc.html"},
		{"README.md", ".pdf", "README.pdf"},
	}
	for _, tt := range tests {
		got := replaceExt(tt.name, tt.newExt)
		if got != tt.want {
			t.Errorf("replaceExt(%q, %q) = %q, want %q", tt.name, tt.newExt, got, tt.want)
		}
	}
}

func TestUnescapeHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"&amp;", "&"},
		{"&lt;div&gt;", "<div>"},
		{"&quot;hello&quot;", `"hello"`},
		{"it&#39;s", "it's"},
		{"&#34;quoted&#34;", `"quoted"`},
		{"no entities", "no entities"},
	}
	for _, tt := range tests {
		got := unescapeHTML(tt.input)
		if got != tt.want {
			t.Errorf("unescapeHTML(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
