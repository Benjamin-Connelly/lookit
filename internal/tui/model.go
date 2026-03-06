package tui

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Benjamin-Connelly/lookit/internal/config"
	"github.com/Benjamin-Connelly/lookit/internal/index"
	"github.com/Benjamin-Connelly/lookit/internal/render"
)

// Panel identifies which panel is currently focused.
type Panel int

const (
	PanelFileList Panel = iota
	PanelPreview
)

// Message types for inter-component communication.
type FileSelectedMsg struct {
	Entry index.FileEntry
}

type PreviewLoadedMsg struct {
	Path    string
	Content string
}

type NavigateMsg struct {
	Path string
}

type StatusMsg struct {
	Text string
}

type clearStatusMsg struct{}

// Model is the root Bubble Tea model for the TUI.
type Model struct {
	cfg      *config.Config
	idx      *index.Index
	links    *index.LinkGraph
	fileList FileListModel
	preview  PreviewModel
	status   StatusBarModel
	keys     KeyMap

	mdRenderer   *render.MarkdownRenderer
	codeRenderer *render.CodeRenderer

	focus    Panel
	width    int
	height   int
	quitting bool
}

// New creates a new root TUI model.
func New(cfg *config.Config, idx *index.Index, links *index.LinkGraph) Model {
	km := DefaultKeyMap()
	switch cfg.Keymap {
	case "vim":
		km = VimKeyMap()
	case "emacs":
		km = EmacsKeyMap()
	}

	mdRenderer, _ := render.NewMarkdownRenderer(cfg.Theme, 80)
	codeRenderer := render.NewCodeRenderer(cfg.Theme, true)

	return Model{
		cfg:          cfg,
		idx:          idx,
		links:        links,
		fileList:     NewFileListModel(idx),
		preview:      NewPreviewModel(),
		status:       NewStatusBarModel(),
		keys:         km,
		mdRenderer:   mdRenderer,
		codeRenderer: codeRenderer,
		focus:        PanelFileList,
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
		if m.fileList.filtering {
			return m.handleFilterKey(msg)
		}
		return m.handleNormalKey(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		previewWidth := m.width - m.width/3 - 1
		m.preview.width = previewWidth
		m.preview.height = m.height - 1
		m.fileList.height = m.height - 1
		if m.mdRenderer != nil {
			_ = m.mdRenderer.SetWidth(previewWidth - 2)
		}
		return m, nil

	case FileSelectedMsg:
		return m.loadPreview(msg.Entry)

	case PreviewLoadedMsg:
		m.preview.SetContent(msg.Path, msg.Content)
		m.status.SetFile(msg.Path)
		return m, nil

	case StatusMsg:
		m.status.SetMessage(msg.Text)
		return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg {
			return clearStatusMsg{}
		})

	case clearStatusMsg:
		m.status.SetMessage("")
		return m, nil
	}

	return m, nil
}

func (m Model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.status.SetMode(m.modeString())
		return m, nil

	case "/", "ctrl+k":
		m.fileList.StartFilter()
		m.status.SetMode("FILTER")
		return m, nil

	case "?":
		content := Help(m.keys)
		m.preview.SetContent("", content)
		m.status.SetFile("Key Bindings")
		return m, nil
	}

	if m.focus == PanelFileList {
		return m.handleFileListKey(msg)
	}
	return m.handlePreviewKey(msg)
}

func (m Model) handleFileListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.fileList.MoveUp()
		return m, nil
	case "down", "j":
		m.fileList.MoveDown()
		return m, nil
	case "enter", "l":
		sel := m.fileList.Selected()
		if sel == nil {
			return m, nil
		}
		return m, func() tea.Msg {
			return FileSelectedMsg{Entry: *sel}
		}
	case "g":
		m.fileList.cursor = 0
		m.fileList.offset = 0
		return m, nil
	case "G":
		if len(m.fileList.filtered) > 0 {
			m.fileList.cursor = len(m.fileList.filtered) - 1
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handlePreviewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.preview.ScrollUp(1)
		return m, nil
	case "down", "j":
		m.preview.ScrollDown(1)
		return m, nil
	case "pgup", "ctrl+u":
		m.preview.ScrollUp(m.preview.height / 2)
		return m, nil
	case "pgdown", "ctrl+d":
		m.preview.ScrollDown(m.preview.height / 2)
		return m, nil
	case "home", "g":
		m.preview.scroll = 0
		return m, nil
	case "end", "G":
		m.preview.ScrollToBottom()
		return m, nil
	}
	return m, nil
}

func (m Model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.fileList.ClearFilter()
		m.status.SetMode("NORMAL")
		return m, nil
	case "enter":
		m.fileList.filtering = false
		m.status.SetMode("NORMAL")
		sel := m.fileList.Selected()
		if sel == nil {
			return m, nil
		}
		return m, func() tea.Msg {
			return FileSelectedMsg{Entry: *sel}
		}
	case "backspace":
		if len(m.fileList.filter) > 0 {
			m.fileList.SetFilter(m.fileList.filter[:len(m.fileList.filter)-1])
		}
		return m, nil
	case "up", "ctrl+p":
		m.fileList.MoveUp()
		return m, nil
	case "down", "ctrl+n":
		m.fileList.MoveDown()
		return m, nil
	default:
		// Single printable characters get appended to the filter
		if len(msg.String()) == 1 {
			m.fileList.SetFilter(m.fileList.filter + msg.String())
		}
		return m, nil
	}
}

func (m Model) loadPreview(entry index.FileEntry) (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		if entry.IsDir {
			return PreviewLoadedMsg{
				Path:    entry.RelPath,
				Content: "[directory]",
			}
		}

		data, err := os.ReadFile(entry.Path)
		if err != nil {
			return PreviewLoadedMsg{
				Path:    entry.RelPath,
				Content: "Error: " + err.Error(),
			}
		}

		content := string(data)

		if entry.IsMarkdown {
			if m.mdRenderer != nil {
				rendered, renderErr := m.mdRenderer.Render(content)
				if renderErr == nil {
					content = rendered
				}
			}
		} else {
			ext := filepath.Ext(entry.RelPath)
			if isTextFile(ext) {
				highlighted, hlErr := m.codeRenderer.Highlight(filepath.Base(entry.RelPath), content)
				if hlErr == nil {
					content = highlighted
				}
			}
		}

		return PreviewLoadedMsg{
			Path:    entry.RelPath,
			Content: content,
		}
	}
}

func (m Model) modeString() string {
	if m.focus == PanelFileList {
		return "FILES"
	}
	return "PREVIEW"
}

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.width == 0 {
		return "Loading..."
	}

	listWidth := m.width / 3
	previewWidth := m.width - listWidth - 1
	contentHeight := m.height - 1

	listStyle := lipgloss.NewStyle().
		Width(listWidth).
		Height(contentHeight)

	previewStyle := lipgloss.NewStyle().
		Width(previewWidth).
		Height(contentHeight)

	// Focus indicator
	var listBorder, previewBorder lipgloss.Style
	if m.focus == PanelFileList {
		listBorder = listStyle.BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63"))
		previewBorder = previewStyle
	} else {
		listBorder = listStyle.BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
		previewBorder = previewStyle
	}

	left := listBorder.Render(m.fileList.View())
	right := previewBorder.Render(m.preview.View())

	main := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	statusWidth := m.width
	m.status.width = statusWidth
	return lipgloss.JoinVertical(lipgloss.Left, main, m.status.View())
}

func isTextFile(ext string) bool {
	ext = strings.ToLower(ext)
	textExts := map[string]bool{
		".go": true, ".js": true, ".ts": true, ".py": true, ".rb": true,
		".rs": true, ".c": true, ".h": true, ".cpp": true, ".java": true,
		".sh": true, ".bash": true, ".zsh": true, ".fish": true,
		".yaml": true, ".yml": true, ".toml": true, ".json": true,
		".xml": true, ".html": true, ".css": true, ".scss": true,
		".sql": true, ".lua": true, ".vim": true, ".el": true,
		".txt": true, ".cfg": true, ".ini": true, ".conf": true,
		".mk": true, ".cmake": true, ".dockerfile": true,
		".gitignore": true, ".env": true, ".mod": true, ".sum": true,
	}
	return textExts[ext]
}
