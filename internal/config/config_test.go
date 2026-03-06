package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Theme != "auto" {
		t.Errorf("default theme should be auto, got %q", cfg.Theme)
	}
	if cfg.Keymap != "default" {
		t.Errorf("default keymap should be default, got %q", cfg.Keymap)
	}
	if cfg.Server.Port != 7777 {
		t.Errorf("default port should be 7777, got %d", cfg.Server.Port)
	}
	if !cfg.Git.Enabled {
		t.Error("git should be enabled by default")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		theme   string
		keymap  string
		wantErr bool
	}{
		{"valid defaults", "auto", "default", false},
		{"dark theme", "dark", "vim", false},
		{"light emacs", "light", "emacs", false},
		{"bad theme", "neon", "default", true},
		{"bad keymap", "auto", "dvorak", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.Theme = tt.theme
			cfg.Keymap = tt.keymap
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadMissing(t *testing.T) {
	// Loading from a nonexistent dir should return defaults
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		// It's acceptable to get defaults or an error for bad paths
		return
	}
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
	if cfg.Theme != "auto" {
		t.Errorf("expected default theme, got %q", cfg.Theme)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	os.WriteFile(cfgPath, []byte(`
theme: dark
keymap: vim
server:
  port: 8080
`), 0o644)

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Theme != "dark" {
		t.Errorf("expected dark theme, got %q", cfg.Theme)
	}
	if cfg.Keymap != "vim" {
		t.Errorf("expected vim keymap, got %q", cfg.Keymap)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}
}

func TestLoadInvalid(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	os.WriteFile(cfgPath, []byte(`
theme: neon
keymap: default
`), 0o644)

	_, err := Load(cfgPath)
	if err == nil {
		t.Error("expected validation error for invalid theme")
	}
}

func TestString(t *testing.T) {
	cfg := DefaultConfig()
	s := cfg.String()
	if s == "" {
		t.Error("String() should not be empty")
	}
}
