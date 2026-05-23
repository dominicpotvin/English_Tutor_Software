package api

import (
	"net/http"

	"englishtutor/internal/model"
)

func (s *Server) handleListQuizzes(w http.ResponseWriter, r *http.Request) {
	quizzes, err := s.store.ListQuizSummaries(r.Context())
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, quizzes)
}

func (s *Server) handleGetQuiz(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx := r.Context()
	quiz, err := s.store.GetQuiz(ctx, id)
	if err != nil {
		s.fail(w, err)
		return
	}
	exercises, err := s.store.ListExercisesByQuiz(ctx, id)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"quiz":      quiz,
		"exercises": exercises,
	})
}

func (s *Server) handleCreateQuiz(w http.ResponseWriter, r *http.Request) {
	var in model.Quiz
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if in.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	quiz, err := s.store.CreateQuiz(r.Context(), in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, quiz)
}

func (s *Server) handleUpdateQuiz(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in model.Quiz
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	quiz, err := s.store.UpdateQuiz(r.Context(), id, in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, quiz)
}

func (s *Server) handleDeleteQuiz(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.store.DeleteQuiz(r.Context(), id); err != nil {
		s.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
