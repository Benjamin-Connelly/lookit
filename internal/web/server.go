package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Benjamin-Connelly/lookit/internal/config"
	"github.com/Benjamin-Connelly/lookit/internal/index"
	"github.com/Benjamin-Connelly/lookit/internal/render"
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

// Start begins listening on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("lookit web server listening on http://%s", addr)
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/", s.handleRoot)
	s.mux.HandleFunc("/__api/files", s.handleAPIFiles)
	s.mux.HandleFunc("/__api/search", s.handleAPISearch)
	s.mux.HandleFunc("/__events", s.handleSSE)
}
