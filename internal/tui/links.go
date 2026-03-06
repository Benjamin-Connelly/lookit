package tui

import (
	"github.com/Benjamin-Connelly/lookit/internal/index"
)

// HistoryEntry records a navigation event for back/forward.
type HistoryEntry struct {
	Path   string
	Scroll int
}

// LinkNavigator manages link following with history.
type LinkNavigator struct {
	graph   *index.LinkGraph
	history []HistoryEntry
	pos     int // current position in history
}

// NewLinkNavigator creates a link navigator backed by a link graph.
func NewLinkNavigator(graph *index.LinkGraph) *LinkNavigator {
	return &LinkNavigator{
		graph: graph,
		pos:   -1,
	}
}

// Navigate pushes a new entry onto the history stack.
func (n *LinkNavigator) Navigate(path string, scroll int) {
	// Truncate forward history
	if n.pos < len(n.history)-1 {
		n.history = n.history[:n.pos+1]
	}
	n.history = append(n.history, HistoryEntry{Path: path, Scroll: scroll})
	n.pos = len(n.history) - 1
}

// Back returns the previous history entry, or nil if at the beginning.
func (n *LinkNavigator) Back() *HistoryEntry {
	if n.pos <= 0 {
		return nil
	}
	n.pos--
	return &n.history[n.pos]
}

// Forward returns the next history entry, or nil if at the end.
func (n *LinkNavigator) Forward() *HistoryEntry {
	if n.pos >= len(n.history)-1 {
		return nil
	}
	n.pos++
	return &n.history[n.pos]
}

// Current returns the current history entry, or nil if empty.
func (n *LinkNavigator) Current() *HistoryEntry {
	if n.pos < 0 || n.pos >= len(n.history) {
		return nil
	}
	return &n.history[n.pos]
}

// LinksAt returns the forward links from the current file.
func (n *LinkNavigator) LinksAt(path string) []index.Link {
	return n.graph.ForwardLinks(path)
}

// BacklinksAt returns files linking to the given path.
func (n *LinkNavigator) BacklinksAt(path string) []index.Link {
	return n.graph.Backlinks(path)
}
