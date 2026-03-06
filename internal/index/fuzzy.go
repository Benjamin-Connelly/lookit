package index

import (
	"github.com/sahilm/fuzzy"
)

// searchSource implements fuzzy.Source for file entries.
type searchSource []FileEntry

func (s searchSource) String(i int) string {
	return s[i].RelPath
}

func (s searchSource) Len() int {
	return len(s)
}

// FuzzySearch performs fuzzy matching against indexed file paths.
// Returns matched entries sorted by relevance.
func (idx *Index) FuzzySearch(query string) []FileEntry {
	idx.mu.RLock()
	entries := make([]FileEntry, len(idx.entries))
	copy(entries, idx.entries)
	idx.mu.RUnlock()

	if query == "" {
		return entries
	}

	matches := fuzzy.FindFrom(query, searchSource(entries))
	result := make([]FileEntry, len(matches))
	for i, m := range matches {
		result[i] = entries[m.Index]
	}
	return result
}
