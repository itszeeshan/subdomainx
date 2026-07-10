package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server is the REST API server for SubdomainX.
type Server struct {
	port      int
	apiKey    string
	outputDir string
	store     *Store
	startTime time.Time
}

// New creates a new API server.
func New(port int, apiKey, outputDir string) *Server {
	return &Server{
		port:      port,
		apiKey:    apiKey,
		outputDir: outputDir,
		store:     NewStore(),
		startTime: time.Now(),
	}
}

// Start registers routes and starts listening. Blocks until the server exits.
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Health endpoint (unauthenticated)
	mux.HandleFunc("GET /api/health", s.handleHealth)

	// Scan endpoints (authenticated)
	mux.HandleFunc("POST /api/scan", s.authMiddleware(s.handleCreateScan))
	mux.HandleFunc("GET /api/scan/{id}", s.authMiddleware(s.handleGetScan))
	mux.HandleFunc("GET /api/scans", s.authMiddleware(s.handleListScans))
	mux.HandleFunc("DELETE /api/scan/{id}", s.authMiddleware(s.handleDeleteScan))

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("SubdomainX API server starting on %s", addr)
	if s.apiKey != "" {
		log.Printf("API key authentication enabled")
	} else {
		log.Printf("WARNING: No API key configured — server is unauthenticated")
	}

	return http.ListenAndServe(addr, mux)
}

// authMiddleware checks the Authorization: Bearer <key> header.
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.apiKey == "" {
			next(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		expected := "Bearer " + s.apiKey
		if auth != expected {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Details: "invalid or missing API key"})
			return
		}
		next(w, r)
	}
}
