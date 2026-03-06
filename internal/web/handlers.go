package web

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	// Prevent path traversal
	cleanPath := filepath.Clean(r.URL.Path)
	if strings.Contains(cleanPath, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	relPath := strings.TrimPrefix(cleanPath, "/")
	if relPath == "" {
		relPath = "."
	}

	entry := s.idx.Lookup(relPath)
	if entry == nil && relPath != "." {
		http.NotFound(w, r)
		return
	}

	if entry != nil && entry.IsDir {
		s.handleDirectory(w, r, relPath)
		return
	}

	if entry != nil && entry.IsMarkdown {
		s.handleMarkdown(w, r, relPath)
		return
	}

	if entry != nil {
		s.handleFile(w, r, relPath)
		return
	}

	// Root directory
	s.handleDirectory(w, r, ".")
}

func (s *Server) handleDirectory(w http.ResponseWriter, r *http.Request, relPath string) {
	entries := s.idx.Entries()
	var dirEntries []map[string]interface{}
	for _, e := range entries {
		dir := filepath.Dir(e.RelPath)
		if dir == "." {
			dir = ""
		}
		target := relPath
		if target == "." {
			target = ""
		}
		if dir == target && e.RelPath != "." {
			dirEntries = append(dirEntries, map[string]interface{}{
				"name":  filepath.Base(e.RelPath),
				"path":  e.RelPath,
				"isDir": e.IsDir,
				"size":  e.Size,
			})
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Placeholder: full template rendering will be implemented
	w.Write([]byte("<html><body><h1>Directory: " + relPath + "</h1>"))
	for _, e := range dirEntries {
		name := e["name"].(string)
		path := e["path"].(string)
		w.Write([]byte("<div><a href=\"/" + path + "\">" + name + "</a></div>"))
	}
	w.Write([]byte("</body></html>"))
}

func (s *Server) handleMarkdown(w http.ResponseWriter, r *http.Request, relPath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Placeholder: Goldmark rendering will be implemented
	w.Write([]byte("<html><body><h1>Markdown: " + relPath + "</h1><p>Rendering not yet implemented</p></body></html>"))
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request, relPath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Placeholder: syntax highlighting will be implemented
	w.Write([]byte("<html><body><h1>File: " + relPath + "</h1><p>Code view not yet implemented</p></body></html>"))
}

func (s *Server) handleAPIFiles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var entries interface{}
	if query != "" {
		entries = s.idx.FuzzySearch(query)
	} else {
		entries = s.idx.Entries()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (s *Server) handleAPISearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "not implemented"})
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	msgCh := make(chan string, 8)
	s.sse.register <- msgCh
	defer func() {
		s.sse.unregister <- msgCh
	}()

	ctx := r.Context()
	for {
		select {
		case msg := <-msgCh:
			w.Write([]byte("data: " + msg + "\n\n"))
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}
