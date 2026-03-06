package index

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors the filesystem for changes and re-indexes affected files.
type Watcher struct {
	index   *Index
	watcher *fsnotify.Watcher
	onChange func(path string) // callback for file changes
	done    chan struct{}
}

// NewWatcher creates a file watcher that updates the index on changes.
func NewWatcher(idx *Index, onChange func(path string)) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		index:   idx,
		watcher: fw,
		onChange: onChange,
		done:    make(chan struct{}),
	}

	go w.loop()

	return w, nil
}

// Start begins watching the index root directory recursively.
func (w *Watcher) Start() error {
	return w.addRecursive(w.index.Root())
}

// Close stops the watcher.
func (w *Watcher) Close() error {
	close(w.done)
	return w.watcher.Close()
}

func (w *Watcher) loop() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				rel, err := filepath.Rel(w.index.Root(), event.Name)
				if err != nil {
					continue
				}
				// Re-build index (could be optimized to incremental)
				if err := w.index.Build(); err != nil {
					log.Printf("watcher: re-index failed: %v", err)
				}
				if w.onChange != nil {
					w.onChange(rel)
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("watcher error: %v", err)
		case <-w.done:
			return
		}
	}
}

func (w *Watcher) addRecursive(root string) error {
	entries := w.index.Entries()
	dirs := make(map[string]bool)
	dirs[root] = true
	for _, e := range entries {
		if e.IsDir && !strings.HasPrefix(filepath.Base(e.Path), ".") {
			dirs[e.Path] = true
		}
	}
	for dir := range dirs {
		if err := w.watcher.Add(dir); err != nil {
			return err
		}
	}
	return nil
}
