package render

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	htmlfmt "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// CodeRenderer provides syntax highlighting via Chroma.
type CodeRenderer struct {
	theme    string
	terminal bool // true for TUI output, false for HTML
}

// NewCodeRenderer creates a code renderer.
// Set terminal=true for TUI (terminal256) output, false for HTML output.
func NewCodeRenderer(theme string, terminal bool) *CodeRenderer {
	return &CodeRenderer{
		theme:    theme,
		terminal: terminal,
	}
}

// Highlight returns syntax-highlighted content for the given file.
func (r *CodeRenderer) Highlight(filename, source string) (string, error) {
	lexer := r.detectLexer(filename, source)
	style := r.getStyle()
	formatter := r.getFormatter()

	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		return source, nil // Fall back to plain text
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return source, nil
	}

	return buf.String(), nil
}

// HighlightToWriter writes syntax-highlighted content to a writer.
func (r *CodeRenderer) HighlightToWriter(w io.Writer, filename, source string) error {
	lexer := r.detectLexer(filename, source)
	style := r.getStyle()
	formatter := r.getFormatter()

	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		_, writeErr := w.Write([]byte(source))
		return writeErr
	}

	return formatter.Format(w, style, iterator)
}

func (r *CodeRenderer) detectLexer(filename, source string) chroma.Lexer {
	ext := strings.ToLower(filepath.Ext(filename))
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Analyse(source)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	_ = ext
	return chroma.Coalesce(lexer)
}

func (r *CodeRenderer) getStyle() *chroma.Style {
	name := "monokai"
	if r.theme == "light" {
		name = "github"
	}
	style := styles.Get(name)
	if style == nil {
		style = styles.Fallback
	}
	return style
}

func (r *CodeRenderer) getFormatter() chroma.Formatter {
	if r.terminal {
		return formatters.Get("terminal256")
	}
	return htmlfmt.New(
		htmlfmt.WithClasses(true),
		htmlfmt.WithLineNumbers(true),
	)
}
