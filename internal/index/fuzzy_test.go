package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFuzzySearch(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("#"), 0o644)
	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "guide.md"), []byte("#"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644)

	idx := New(dir)
	idx.Build()

	// Empty query returns all
	all := idx.FuzzySearch("")
	if len(all) == 0 {
		t.Error("empty query should return all entries")
	}

	// Search for "read" should match README.md
	results := idx.FuzzySearch("read")
	found := false
	for _, r := range results {
		if r.RelPath == "README.md" {
			found = true
			break
		}
	}
	if !found {
		t.Error("fuzzy search for 'read' should find README.md")
	}

	// maxResults limit
	limited := idx.FuzzySearch("", 1)
	if len(limited) != 1 {
		t.Errorf("expected 1 result with limit, got %d", len(limited))
	}
}

func TestFuzzySearchMarkdown(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("#"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644)

	idx := New(dir)
	idx.Build()

	results := idx.FuzzySearchMarkdown("")
	for _, r := range results {
		if !r.IsMarkdown {
			t.Errorf("FuzzySearchMarkdown returned non-markdown: %s", r.RelPath)
		}
	}
}
