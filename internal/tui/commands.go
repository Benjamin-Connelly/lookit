package tui

// CommandEntry represents an entry in the command palette.
type CommandEntry struct {
	Name        string
	Description string
	Action      func() // executed when selected
}

// CommandPalette manages the colon-mode command interface.
type CommandPalette struct {
	commands []CommandEntry
	input    string
	filtered []CommandEntry
	cursor   int
	active   bool
}

// NewCommandPalette creates a command palette with default commands.
func NewCommandPalette() CommandPalette {
	return CommandPalette{}
}

// RegisterCommand adds a command to the palette.
func (p *CommandPalette) RegisterCommand(cmd CommandEntry) {
	p.commands = append(p.commands, cmd)
}

// Open activates the command palette.
func (p *CommandPalette) Open() {
	p.active = true
	p.input = ""
	p.filtered = p.commands
	p.cursor = 0
}

// Close deactivates the command palette.
func (p *CommandPalette) Close() {
	p.active = false
	p.input = ""
}

// IsActive returns whether the palette is open.
func (p *CommandPalette) IsActive() bool {
	return p.active
}

// SetInput updates the filter input.
func (p *CommandPalette) SetInput(s string) {
	p.input = s
	// Simple prefix match for now
	p.filtered = nil
	for _, cmd := range p.commands {
		if len(p.input) == 0 || containsIgnoreCase(cmd.Name, p.input) {
			p.filtered = append(p.filtered, cmd)
		}
	}
	p.cursor = 0
}

// Execute runs the selected command.
func (p *CommandPalette) Execute() {
	if p.cursor >= 0 && p.cursor < len(p.filtered) {
		if p.filtered[p.cursor].Action != nil {
			p.filtered[p.cursor].Action()
		}
	}
	p.Close()
}

// View renders the command palette overlay.
func (p CommandPalette) View() string {
	if !p.active {
		return ""
	}
	s := ":" + p.input + "\n"
	for i, cmd := range p.filtered {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		s += cursor + cmd.Name + " - " + cmd.Description + "\n"
	}
	return s
}

func containsIgnoreCase(s, substr string) bool {
	// Simple case-insensitive contains
	sl := len(s)
	subl := len(substr)
	if subl > sl {
		return false
	}
	for i := 0; i <= sl-subl; i++ {
		match := true
		for j := 0; j < subl; j++ {
			a, b := s[i+j], substr[j]
			if a >= 'A' && a <= 'Z' {
				a += 32
			}
			if b >= 'A' && b <= 'Z' {
				b += 32
			}
			if a != b {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
