package web

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Benjamin-Connelly/lookit/internal/config"
	"github.com/Benjamin-Connelly/lookit/internal/index"
	"github.com/Benjamin-Connelly/lookit/internal/render"
	"github.com/Benjamin-Connelly/lookit/internal/web/static"
)

// Server is the HTTP server for web mode.
type Server struct {
	cfg    *config.Config
	idx    *index.Index
	links  *index.LinkGraph
	code   *render.CodeRenderer
	mux    *http.ServeMux
	server *http.Server
	sse    *SSEBroker
}

// SSEBroker manages Server-Sent Events for live reload.
type SSEBroker struct {
	clients    map[chan string]bool
	register   chan chan string
	unregister chan chan string
	broadcast  chan string
}

// NewSSEBroker creates a new SSE event broker.
func NewSSEBroker() *SSEBroker {
	b := &SSEBroker{
		clients:    make(map[chan string]bool),
		register:   make(chan chan string),
		unregister: make(chan chan string),
		broadcast:  make(chan string, 16),
	}
	go b.run()
	return b
}

func (b *SSEBroker) run() {
	for {
		select {
		case client := <-b.register:
			b.clients[client] = true
		case client := <-b.unregister:
			delete(b.clients, client)
			close(client)
		case msg := <-b.broadcast:
			for client := range b.clients {
				select {
				case client <- msg:
				default:
					delete(b.clients, client)
					close(client)
				}
			}
		}
	}
}

// Notify sends a reload event to all connected clients.
func (b *SSEBroker) Notify(path string) {
	b.broadcast <- path
}

// New creates a new web server.
func New(cfg *config.Config, idx *index.Index, links *index.LinkGraph) *Server {
	s := &Server{
		cfg:   cfg,
		idx:   idx,
		links: links,
		code:  render.NewCodeRenderer(cfg.Theme, false),
		mux:   http.NewServeMux(),
		sse:   NewSSEBroker(),
	}

	s.registerRoutes()
	return s
}

// OnFileChange is a callback for the file watcher. Wire it to index.Watcher's onChange.
func (s *Server) OnFileChange(relPath string) {
	s.sse.Notify(relPath)
}

// Start begins listening on the configured port and handles graceful shutdown.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.middleware(s.mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown on SIGINT/SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		log.Printf("lookit web server listening on http://%s", addr)
		errCh <- s.server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case sig := <-sigCh:
		log.Printf("received %v, shutting down", sig)
		return s.Stop()
	}
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// middleware chains security headers, request logging, and ETag support.
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'")

		// ETag support for non-SSE, non-API requests
		if !strings.HasPrefix(r.URL.Path, "/__events") && !strings.HasPrefix(r.URL.Path, "/__api/") {
			rec := &responseRecorder{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(rec, r)

			// Log request
			log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.statusCode, time.Since(start).Round(time.Millisecond))
			return
		}

		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}

// etagMiddleware wraps a handler to add ETag caching.
func etagMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{ResponseWriter: w, statusCode: 200}
		next.ServeHTTP(rec, r)

		if rec.statusCode == 200 && len(rec.body) > 0 {
			etag := fmt.Sprintf(`"%x"`, md5.Sum(rec.body))
			w.Header().Set("ETag", etag)
			w.Header().Set("Cache-Control", "no-cache")

			if r.Header.Get("If-None-Match") == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			w.Header().Set("Content-Type", rec.contentType)
			w.WriteHeader(rec.statusCode)
			w.Write(rec.body)
			return
		}

		// Non-200 or empty body: already written by recorder fallthrough
		if !rec.captured {
			return
		}
		w.Header().Set("Content-Type", rec.contentType)
		w.WriteHeader(rec.statusCode)
		w.Write(rec.body)
	}
}

// responseRecorder captures response data for ETag generation.
type responseRecorder struct {
	http.ResponseWriter
	statusCode  int
	body        []byte
	contentType string
	captured    bool
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.captured = true
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.captured = true
	r.contentType = r.Header().Get("Content-Type")
	r.body = append(r.body, b...)
	return len(b), nil
}

func (s *Server) registerRoutes() {
	// Static assets
	staticFS, err := fs.Sub(static.Files, ".")
	if err != nil {
		log.Fatalf("failed to create static sub-filesystem: %v", err)
	}
	s.mux.Handle("/__static/", http.StripPrefix("/__static/", http.FileServer(http.FS(staticFS))))

	// API routes
	s.mux.HandleFunc("/__api/files", s.handleAPIFiles)
	s.mux.HandleFunc("/__api/search", s.handleAPISearch)
	s.mux.HandleFunc("/__events", s.handleSSE)

	// All other routes go through root handler with ETag support
	s.mux.HandleFunc("/", etagMiddleware(s.handleRoot))
}
