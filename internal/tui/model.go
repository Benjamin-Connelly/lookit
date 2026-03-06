package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Benjamin-Connelly/lookit/internal/config"
	"github.com/Benjamin-Connelly/lookit/internal/index"
)

// Panel identifies which panel is currently focused.
type Panel int

const (
	PanelFileList Panel = iota
	PanelPreview
)

// Model is the root Bubble Tea model for the TUI.
type Model struct {
	cfg      *config.Config
	idx      *index.Index
	links    *index.LinkGraph
	fileList FileListModel
	preview  PreviewModel
	status   StatusBarModel
	keys     KeyMap

	focus    Panel
	width    int
	height   int
	quitting bool
}

// New creates a new root TUI model.
func New(cfg *config.Config, idx *index.Index, links *index.LinkGraph) Model {
	return Model{
		cfg:      cfg,
		idx:      idx,
		links:    links,
		fileList: NewFileListModel(idx),
		preview:  NewPreviewModel(),
		status:   NewStatusBarModel(),
		keys:     DefaultKeyMap(),
		focus:    PanelFileList,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			if m.focus == PanelFileList {
				m.focus = PanelPreview
			} else {
				m.focus = PanelFileList
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	listWidth := m.width / 3
	previewWidth := m.width - listWidth - 1

	listStyle := lipgloss.NewStyle().Width(listWidth).Height(m.height - 1)
	previewStyle := lipgloss.NewStyle().Width(previewWidth).Height(m.height - 1)

	left := listStyle.Render(m.fileList.View())
	right := previewStyle.Render(m.preview.View())

	main := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)
	return lipgloss.JoinVertical(lipgloss.Left, main, m.status.View())
}
