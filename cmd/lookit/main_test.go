package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Benjamin-Connelly/lookit/internal/config"
)

func init() {
	cfg = config.DefaultConfig()
}

func TestResolveRoot_Directory(t *testing.T) {
	dir := t.TempDir()

	root, initialFile, err := resolveRoot([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root != dir {
		t.Errorf("root = %q, want %q", root, dir)
	}
	if initialFile != "" {
		t.Errorf("initialFile = %q, want empty", initialFile)
	}
}

func TestResolveRoot_SingleFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "README.md")
	os.WriteFile(file, []byte("# Hello\n"), 0o644)

	root, initialFile, err := resolveRoot([]string{file})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root != dir {
		t.Errorf("root = %q, want %q", root, dir)
	}
	if initialFile != "README.md" {
		t.Errorf("initialFile = %q, want %q", initialFile, "README.md")
	}
}

func TestResolveRoot_NonExistent(t *testing.T) {
	_, _, err := resolveRoot([]string{"/nonexistent/path"})
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}

func TestResolveRoot_DefaultConfig(t *testing.T) {
	oldRoot := cfg.Root
	cfg.Root = t.TempDir()
	defer func() { cfg.Root = oldRoot }()

	root, initialFile, err := resolveRoot(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root != cfg.Root {
		t.Errorf("root = %q, want %q", root, cfg.Root)
	}
	if initialFile != "" {
		t.Errorf("initialFile = %q, want empty", initialFile)
	}
}
