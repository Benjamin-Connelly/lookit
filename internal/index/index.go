package index

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FileEntry represents an indexed file.
type FileEntry struct {
	Path    string
	RelPath string
	Size    int64
	ModTime time.Time
	IsDir   bool
	IsMarkdown bool
}

// Index maintains an in-memory file index for fast lookup and search.
type Index struct {
	root    string
	entries []FileEntry
	byPath  map[string]*FileEntry
	mu      sync.RWMutex
}

// New creates a new Index rooted at the given directory.
func New(root string) *Index {
	return &Index{
		root:   root,
		byPath: make(map[string]*FileEntry),
	}
}

// Build walks the root directory and populates the index.
func (idx *Index) Build() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.entries = nil
	idx.byPath = make(map[string]*FileEntry)

	return filepath.WalkDir(idx.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip files we can't read
		}

		// Skip hidden directories
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && path != idx.root {
			return filepath.SkipDir
		}

		rel, err := filepath.Rel(idx.root, path)
		if err != nil {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		entry := FileEntry{
			Path:       path,
			RelPath:    rel,
			Size:       info.Size(),
			ModTime:    info.ModTime(),
			IsDir:      d.IsDir(),
			IsMarkdown: isMarkdown(d.Name()),
		}

		idx.entries = append(idx.entries, entry)
		idx.byPath[rel] = &idx.entries[len(idx.entries)-1]

		return nil
	})
}

// Entries returns a copy of all indexed entries.
func (idx *Index) Entries() []FileEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	result := make([]FileEntry, len(idx.entries))
	copy(result, idx.entries)
	return result
}

// MarkdownFiles returns only markdown file entries.
func (idx *Index) MarkdownFiles() []FileEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	var result []FileEntry
	for _, e := range idx.entries {
		if e.IsMarkdown {
			result = append(result, e)
		}
	}
	return result
}

// Lookup returns the entry for a relative path, or nil if not found.
func (idx *Index) Lookup(relPath string) *FileEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.byPath[relPath]
}

// Root returns the index root directory.
func (idx *Index) Root() string {
	return idx.root
}

func isMarkdown(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".md" || ext == ".markdown" || ext == ".mdown"
}
