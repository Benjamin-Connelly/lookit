package tui

// PreviewModel renders file content in the preview pane.
type PreviewModel struct {
	content  string
	filePath string
	scroll   int
}

// NewPreviewModel creates a preview pane.
func NewPreviewModel() PreviewModel {
	return PreviewModel{}
}

// SetContent updates the preview with rendered content.
func (m *PreviewModel) SetContent(path, content string) {
	m.filePath = path
	m.content = content
	m.scroll = 0
}

// ScrollUp scrolls the preview up.
func (m *PreviewModel) ScrollUp(lines int) {
	m.scroll -= lines
	if m.scroll < 0 {
		m.scroll = 0
	}
}

// ScrollDown scrolls the preview down.
func (m *PreviewModel) ScrollDown(lines int) {
	m.scroll += lines
}

// View renders the preview content.
func (m PreviewModel) View() string {
	if m.content == "" {
		return "Select a file to preview"
	}
	return m.content
}
