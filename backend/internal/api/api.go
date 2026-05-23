// Package api exposes the HTTP REST interface over the store.
package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"englishtutor/internal/store"
)

// Server holds the dependencies shared by every HTTP handler.
type Server struct {
	store *store.Store
}

// NewServer builds a Server backed by the given store.
func NewServer(st *store.Store) *Server {
	return &Server{store: st}
}

// Handler builds the fully-wired HTTP handler, middleware included.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)

	mux.HandleFunc("GET /api/levels", s.handleListLevels)
	mux.HandleFunc("POST /api/levels", s.handleCreateLevel)
	mux.HandleFunc("GET /api/levels/{id}", s.handleGetLevel)
	mux.HandleFunc("PUT /api/levels/{id}", s.handleUpdateLevel)

	mux.HandleFunc("GET /api/lessons/{id}", s.handleGetLesson)
	mux.HandleFunc("POST /api/lessons", s.handleCreateLesson)
	mux.HandleFunc("PUT /api/lessons/{id}", s.handleUpdateLesson)
	mux.HandleFunc("DELETE /api/lessons/{id}", s.handleDeleteLesson)

	mux.HandleFunc("GET /api/topics/{id}", s.handleGetTopic)
	mux.HandleFunc("GET /api/topics/{id}/exercises", s.handleListTopicExercises)
	mux.HandleFunc("POST /api/topics", s.handleCreateTopic)
	mux.HandleFunc("PUT /api/topics/{id}", s.handleUpdateTopic)
	mux.HandleFunc("DELETE /api/topics/{id}", s.handleDeleteTopic)

	mux.HandleFunc("GET /api/exercises/{id}", s.handleGetExercise)
	mux.HandleFunc("POST /api/exercises", s.handleCreateExercise)
	mux.HandleFunc("PUT /api/exercises/{id}", s.handleUpdateExercise)
	mux.HandleFunc("DELETE /api/exercises/{id}", s.handleDeleteExercise)
	mux.HandleFunc("POST /api/exercises/{id}/attempt", s.handleAttempt)

	mux.HandleFunc("GET /api/quizzes", s.handleListQuizzes)
	mux.HandleFunc("POST /api/quizzes", s.handleCreateQuiz)
	mux.HandleFunc("GET /api/quizzes/{id}", s.handleGetQuiz)
	mux.HandleFunc("PUT /api/quizzes/{id}", s.handleUpdateQuiz)
	mux.HandleFunc("DELETE /api/quizzes/{id}", s.handleDeleteQuiz)

	mux.HandleFunc("GET /api/vocabulary", s.handleListVocabulary)
	mux.HandleFunc("POST /api/vocabulary", s.handleCreateVocabulary)
	mux.HandleFunc("PUT /api/vocabulary/{id}", s.handleUpdateVocabulary)
	mux.HandleFunc("DELETE /api/vocabulary/{id}", s.handleDeleteVocabulary)

	mux.HandleFunc("GET /api/progress", s.handleProgress)

	return logRequests(withCORS(mux))
}

// --- Response helpers -----------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("encode response: %v", err)
		}
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// fail maps a store error to the right HTTP status.
func (s *Server) fail(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}
	log.Printf("internal error: %v", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}

// --- Request helpers ------------------------------------------------------

func pathID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func decodeBody(r *http.Request, dst any) bool {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst) == nil
}

// optionalID reads a positive int64 query parameter, returning nil when absent.
func optionalID(r *http.Request, key string) *int64 {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil
	}
	return &id
}

// optionalString reads a query parameter, returning nil when absent.
func optionalString(r *http.Request, key string) *string {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil
	}
	return &raw
}

// --- Middleware -----------------------------------------------------------

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Millisecond))
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
