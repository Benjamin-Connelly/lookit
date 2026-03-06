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
// Returns matched entries sorted by relevance. An optional maxResults
// limits the number of results (0 or omitted means unlimited).
func (idx *Index) FuzzySearch(query string, maxResults ...int) []FileEntry {
	limit := 0
	if len(maxResults) > 0 {
		limit = maxResults[0]
	}

	idx.mu.RLock()
	entries := make([]FileEntry, len(idx.entries))
	copy(entries, idx.entries)
	idx.mu.RUnlock()

	if query == "" {
		if limit > 0 && limit < len(entries) {
			return entries[:limit]
		}
		return entries
	}

	matches := fuzzy.FindFrom(query, searchSource(entries))
	result := make([]FileEntry, 0, len(matches))
	for _, m := range matches {
		result = append(result, entries[m.Index])
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}

// FuzzySearchMarkdown performs fuzzy matching against only markdown file paths.
// An optional maxResults limits the number of results (0 or omitted means unlimited).
func (idx *Index) FuzzySearchMarkdown(query string, maxResults ...int) []FileEntry {
	limit := 0
	if len(maxResults) > 0 {
		limit = maxResults[0]
	}

	mdFiles := idx.MarkdownFiles()

	if query == "" {
		if limit > 0 && limit < len(mdFiles) {
			return mdFiles[:limit]
		}
		return mdFiles
	}

	matches := fuzzy.FindFrom(query, searchSource(mdFiles))
	result := make([]FileEntry, 0, len(matches))
	for _, m := range matches {
		result = append(result, mdFiles[m.Index])
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}
