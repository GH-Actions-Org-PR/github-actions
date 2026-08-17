// Package handler contains the HTTP handlers for the API.
package handler

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/example/go-api/internal/version"
)

// Server bundles handler state. Kept small and dependency-injected so
// tests can construct it without spinning up real infrastructure.
type Server struct {
	ready   atomic.Bool
	started time.Time
}

func New() *Server {
	s := &Server{started: time.Now()}
	s.ready.Store(true)
	return s
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("GET /readyz", s.handleReadyz)
	mux.HandleFunc("GET /version", s.handleVersion)
	mux.HandleFunc("GET /api/v1/ping", s.handlePing)
	return withLogging(mux)
}

// handleHealthz reports process liveness - used by the container
// orchestrator's liveness probe. It never depends on downstream systems.
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleReadyz reports whether the service is ready to accept traffic -
// used by the readiness probe. SetNotReady is called during graceful
// shutdown so the load balancer stops routing new requests first.
func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	if !s.ready.Load() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "draining"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"uptime": time.Since(s.started).String(),
	})
}

func (s *Server) SetNotReady() { s.ready.Store(false) }

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, version.Info{
		Version:   version.Version,
		Commit:    version.Commit,
		BuildDate: version.BuildDate,
		GoVersion: runtime.Version(),
	})
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		_ = time.Since(start) // wire up to a real logger/metrics client in production
	})
}
