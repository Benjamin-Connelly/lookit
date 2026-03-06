package render

import (
	"bytes"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Heading represents a markdown heading extracted from source.
type Heading struct {
	Level int
	Text  string
	Line  int
}

// Link represents a markdown link extracted from source.
type Link struct {
	Text        string
	Destination string
	Line        int
}

// MarkdownRenderer wraps Glamour for TUI markdown rendering.
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
	theme    string
	width    int
}

// NewMarkdownRenderer creates a markdown renderer with the given theme and width.
func NewMarkdownRenderer(theme string, width int) (*MarkdownRenderer, error) {
	styleName := resolveTheme(theme)

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

// resolveTheme maps theme names to Glamour style names, with auto-detection.
func resolveTheme(theme string) string {
	switch theme {
	case "light":
		return "light"
	case "auto":
		if lipgloss.HasDarkBackground() {
			return "dark"
		}
		return "light"
	default:
		return "dark"
	}
}

// Render converts markdown to styled terminal output.
// On error, returns the raw source as fallback.
func (r *MarkdownRenderer) Render(source string) (string, error) {
	out, err := r.renderer.Render(source)
	if err != nil {
		return source, nil
	}
	return out, nil
}

// RenderFile reads a file and renders its markdown content.
// On render error, returns the raw file content as fallback.
func (r *MarkdownRenderer) RenderFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return r.Render(string(data))
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

// parseMarkdown parses source into a goldmark AST.
func parseMarkdown(source []byte) ast.Node {
	md := goldmark.New()
	reader := text.NewReader(source)
	return md.Parser().Parse(reader)
}

// lineNumber returns the 1-based line number for a byte offset in source.
func lineNumber(source []byte, offset int) int {
	if offset > len(source) {
		offset = len(source)
	}
	return bytes.Count(source[:offset], []byte("\n")) + 1
}

// nodeText extracts the text content of a node from source.
func nodeText(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			buf.Write(t.Segment.Value(source))
		}
	}
	return buf.String()
}

// nodeStartOffset returns the byte offset where a node starts in source.
func nodeStartOffset(n ast.Node) int {
	if n.Type() == ast.TypeBlock {
		if bl, ok := n.(interface{ Lines() *text.Segments }); ok {
			if bl.Lines().Len() > 0 {
				return bl.Lines().At(0).Start
			}
		}
	}
	// For inline nodes, walk children to find first text segment
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			return t.Segment.Start
		}
	}
	return 0
}

// ExtractHeadings returns all headings from markdown source.
func ExtractHeadings(source string) []Heading {
	src := []byte(source)
	doc := parseMarkdown(src)

	var headings []Heading
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if h, ok := n.(*ast.Heading); ok {
			headings = append(headings, Heading{
				Level: h.Level,
				Text:  nodeText(h, src),
				Line:  lineNumber(src, nodeStartOffset(h)),
			})
		}
		return ast.WalkContinue, nil
	})
	return headings
}

// ExtractLinks returns all links from markdown source.
func ExtractLinks(source string) []Link {
	src := []byte(source)
	doc := parseMarkdown(src)

	var links []Link
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if l, ok := n.(*ast.Link); ok {
			links = append(links, Link{
				Text:        nodeText(l, src),
				Destination: string(l.Destination),
				Line:        lineNumber(src, nodeStartOffset(l)),
			})
		}
		return ast.WalkContinue, nil
	})
	return links
}
