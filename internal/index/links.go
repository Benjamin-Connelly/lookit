package index

import (
	"sync"
)

// Link represents a link from one file to another.
type Link struct {
	Source string // relative path of the source file
	Target string // relative path or URL of the target
	Text   string // link text
	Line   int    // line number in source
	Broken bool   // true if target cannot be resolved
}

// LinkGraph maintains bidirectional link relationships between files.
type LinkGraph struct {
	// Forward links: source -> []Link
	forward map[string][]Link
	// Backlinks: target -> []Link
	backward map[string][]Link
	mu       sync.RWMutex
}

// NewLinkGraph creates an empty link graph.
func NewLinkGraph() *LinkGraph {
	return &LinkGraph{
		forward:  make(map[string][]Link),
		backward: make(map[string][]Link),
	}
}

// SetLinks replaces all links originating from the given source file.
func (g *LinkGraph) SetLinks(source string, links []Link) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Remove old backlinks for this source
	for _, oldLink := range g.forward[source] {
		g.removeBacklink(oldLink.Target, source)
	}

	// Set new forward links
	g.forward[source] = links

	// Add new backlinks
	for _, link := range links {
		g.backward[link.Target] = append(g.backward[link.Target], link)
	}
}

// ForwardLinks returns all links originating from the given file.
func (g *LinkGraph) ForwardLinks(source string) []Link {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]Link, len(g.forward[source]))
	copy(result, g.forward[source])
	return result
}

// Backlinks returns all links pointing to the given file.
func (g *LinkGraph) Backlinks(target string) []Link {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]Link, len(g.backward[target]))
	copy(result, g.backward[target])
	return result
}

// BrokenLinks returns all links in the graph that are marked as broken.
func (g *LinkGraph) BrokenLinks() []Link {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var broken []Link
	for _, links := range g.forward {
		for _, link := range links {
			if link.Broken {
				broken = append(broken, link)
			}
		}
	}
	return broken
}

func (g *LinkGraph) removeBacklink(target, source string) {
	backlinks := g.backward[target]
	filtered := backlinks[:0]
	for _, bl := range backlinks {
		if bl.Source != source {
			filtered = append(filtered, bl)
		}
	}
	g.backward[target] = filtered
}
