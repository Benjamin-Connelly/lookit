package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// StatusBarModel renders the bottom status bar.
type StatusBarModel struct {
	filePath string
	message  string
	mode     string
}

// NewStatusBarModel creates a status bar.
func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		mode: "NORMAL",
	}
}

// SetFile updates the displayed file path.
func (m *StatusBarModel) SetFile(path string) {
	m.filePath = path
}

// SetMessage sets a temporary status message.
func (m *StatusBarModel) SetMessage(msg string) {
	m.message = msg
}

// SetMode sets the current input mode display.
func (m *StatusBarModel) SetMode(mode string) {
	m.mode = mode
}

// View renders the status bar.
func (m StatusBarModel) View() string {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252")).
		Padding(0, 1)

	left := m.mode
	if m.filePath != "" {
		left += " | " + m.filePath
	}

	right := ""
	if m.message != "" {
		right = m.message
	}

	return style.Render(fmt.Sprintf("%-40s %s", left, right))
}
