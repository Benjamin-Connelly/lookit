package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const maxRecentFiles = 50

// RecentFiles manages a persistent list of recently opened files.
type RecentFiles struct {
	Files []string `json:"files"`
	path  string
}

// LoadRecentFiles reads the recent files list from disk.
func LoadRecentFiles() *RecentFiles {
	r := &RecentFiles{}
	configDir, err := ConfigDir()
	if err != nil {
		return r
	}
	r.path = filepath.Join(configDir, "recent.json")

	data, err := os.ReadFile(r.path)
	if err != nil {
		return r
	}
	_ = json.Unmarshal(data, r)
	return r
}

// Add puts a file at the front of the recent list, removing duplicates.
func (r *RecentFiles) Add(path string) {
	// Remove existing occurrence
	filtered := make([]string, 0, len(r.Files))
	for _, f := range r.Files {
		if f != path {
			filtered = append(filtered, f)
		}
	}
	// Prepend
	r.Files = append([]string{path}, filtered...)
	if len(r.Files) > maxRecentFiles {
		r.Files = r.Files[:maxRecentFiles]
	}
}

// Save writes the recent files list to disk.
func (r *RecentFiles) Save() error {
	if r.path == "" {
		configDir, err := ConfigDir()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			return err
		}
		r.path = filepath.Join(configDir, "recent.json")
	}
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0o644)
}
