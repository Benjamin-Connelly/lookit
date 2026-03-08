package tui

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ImageRenderer handles inline image display in the TUI.
// Uses terminal-specific protocols (Kitty, iTerm2, Sixel) when available.
type ImageRenderer struct {
	protocol ImageProtocol
	maxWidth int
}

// ImageProtocol identifies the terminal image protocol to use.
type ImageProtocol int

const (
	ImageProtocolNone ImageProtocol = iota
	ImageProtocolKitty
	ImageProtocolITerm2
)

// NewImageRenderer creates an image renderer, auto-detecting protocol.
func NewImageRenderer() *ImageRenderer {
	return &ImageRenderer{
		protocol: detectImageProtocol(),
		maxWidth: 80,
	}
}

// CanRender returns whether inline images are supported.
func (r *ImageRenderer) CanRender() bool {
	return r.protocol != ImageProtocolNone
}

// ProtocolName returns the detected protocol name for display.
func (r *ImageRenderer) ProtocolName() string {
	switch r.protocol {
	case ImageProtocolKitty:
		return "kitty"
	case ImageProtocolITerm2:
		return "iterm2"
	default:
		return "none"
	}
}

// Render returns a string that will display the image inline.
func (r *ImageRenderer) Render(path string) string {
	if !r.CanRender() {
		return "[image: " + filepath.Base(path) + "]"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "[image: error reading " + filepath.Base(path) + "]"
	}

	switch r.protocol {
	case ImageProtocolKitty:
		return r.renderKitty(data)
	case ImageProtocolITerm2:
		return r.renderITerm2(data, path)
	default:
		return "[image: " + filepath.Base(path) + "]"
	}
}

// renderKitty renders an image using the Kitty graphics protocol.
// See: https://sw.kovidgoyal.net/kitty/graphics-protocol/
func (r *ImageRenderer) renderKitty(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)

	// Kitty protocol: split into chunks of 4096 bytes
	var b strings.Builder
	for i := 0; i < len(encoded); i += 4096 {
		end := i + 4096
		if end > len(encoded) {
			end = len(encoded)
		}
		chunk := encoded[i:end]

		more := 1
		if end >= len(encoded) {
			more = 0
		}

		if i == 0 {
			// First chunk: include format and action
			b.WriteString(fmt.Sprintf("\033_Ga=T,f=100,m=%d;%s\033\\", more, chunk))
		} else {
			b.WriteString(fmt.Sprintf("\033_Gm=%d;%s\033\\", more, chunk))
		}
	}
	return b.String()
}

// renderITerm2 renders an image using the iTerm2 inline image protocol.
// Also works in WezTerm and other compatible terminals.
// See: https://iterm2.com/documentation-images.html
func (r *ImageRenderer) renderITerm2(data []byte, path string) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	name := base64.StdEncoding.EncodeToString([]byte(filepath.Base(path)))
	return fmt.Sprintf("\033]1337;File=name=%s;inline=1;width=auto;preserveAspectRatio=1:%s\a",
		name, encoded)
}

func detectImageProtocol() ImageProtocol {
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")

	// Kitty detection
	if strings.Contains(term, "kitty") || termProgram == "kitty" {
		return ImageProtocolKitty
	}

	// Ghostty uses Kitty graphics protocol
	if termProgram == "ghostty" {
		return ImageProtocolKitty
	}

	// iTerm2 and compatible terminals
	if termProgram == "iTerm.app" || termProgram == "WezTerm" {
		return ImageProtocolITerm2
	}

	// Check LC_TERMINAL for iTerm2 (set inside tmux/screen)
	if os.Getenv("LC_TERMINAL") == "iTerm2" {
		return ImageProtocolITerm2
	}

	return ImageProtocolNone
}

// IsImageFile returns whether the given extension is a supported image format.
func IsImageFile(ext string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".svg":
		return true
	}
	return false
}
