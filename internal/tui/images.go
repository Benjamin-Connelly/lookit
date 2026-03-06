package tui

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
	ImageProtocolSixel
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

// Render returns a string that will display the image inline.
func (r *ImageRenderer) Render(path string) string {
	if !r.CanRender() {
		return "[image: " + path + "]"
	}
	// Protocol-specific rendering will be implemented per protocol
	return "[image: " + path + "]"
}

func detectImageProtocol() ImageProtocol {
	// Detection logic for terminal capabilities
	// Will check TERM_PROGRAM, TERM, etc.
	return ImageProtocolNone
}
