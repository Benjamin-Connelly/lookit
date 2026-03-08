package tui

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

// ImageRenderer handles image file display in the TUI.
// Uses a text-based info card because terminal image protocols (Kitty, iTerm2)
// are incompatible with Bubble Tea's alt-screen rendering — they cause
// duplication on resize and can't be scrolled or clipped to pane bounds.
type ImageRenderer struct{}

// NewImageRenderer creates an image renderer.
func NewImageRenderer() *ImageRenderer {
	return &ImageRenderer{}
}

// Render returns a text-based info card for the image file.
func (r *ImageRenderer) Render(path string) string {
	name := filepath.Base(path)
	ext := strings.ToLower(filepath.Ext(path))

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf("  [image: %s — cannot read file]", name)
	}

	size := formatFileSize(info.Size())

	// Try to read image dimensions
	dims := ""
	if f, err := os.Open(path); err == nil {
		defer f.Close()
		if cfg, _, err := image.DecodeConfig(f); err == nil {
			dims = fmt.Sprintf("%d × %d", cfg.Width, cfg.Height)
		}
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString("  ┌─────────────────────────────────────┐\n")
	b.WriteString(fmt.Sprintf("  │  %s\n", name))
	b.WriteString("  │\n")
	b.WriteString(fmt.Sprintf("  │  Type:  %s\n", strings.ToUpper(strings.TrimPrefix(ext, "."))))
	b.WriteString(fmt.Sprintf("  │  Size:  %s\n", size))
	if dims != "" {
		b.WriteString(fmt.Sprintf("  │  Dims:  %s\n", dims))
	}
	b.WriteString("  │\n")
	b.WriteString("  │  Press 'e' to open in external viewer\n")
	b.WriteString("  └─────────────────────────────────────┘\n")

	return b.String()
}

// IsImageFile returns whether the given extension is a supported image format.
func IsImageFile(ext string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".svg", ".ico":
		return true
	}
	return false
}

func formatFileSize(bytes int64) string {
	switch {
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
