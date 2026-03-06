package render

import (
	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer wraps Glamour for TUI markdown rendering.
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
	theme    string
	width    int
}

// NewMarkdownRenderer creates a markdown renderer with the given theme and width.
func NewMarkdownRenderer(theme string, width int) (*MarkdownRenderer, error) {
	styleName := "dark"
	if theme == "light" {
		styleName = "light"
	} else if theme == "auto" {
		styleName = "auto"
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(styleName),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	return &MarkdownRenderer{
		renderer: r,
		theme:    theme,
		width:    width,
	}, nil
}

// Render converts markdown to styled terminal output.
func (r *MarkdownRenderer) Render(source string) (string, error) {
	return r.renderer.Render(source)
}

// SetWidth updates the word wrap width and recreates the renderer.
func (r *MarkdownRenderer) SetWidth(width int) error {
	r.width = width
	nr, err := NewMarkdownRenderer(r.theme, width)
	if err != nil {
		return err
	}
	r.renderer = nr.renderer
	return nil
}
