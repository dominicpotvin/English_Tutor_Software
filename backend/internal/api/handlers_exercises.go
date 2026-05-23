package api

import (
	"net/http"

	"englishtutor/internal/model"
)

// validateExercise returns an error message, or "" when the exercise is valid.
func validateExercise(e model.Exercise) string {
	switch e.Kind {
	case "mcq", "fill_blank", "true_false":
	default:
		return "kind must be mcq, fill_blank or true_false"
	}
	if (e.TopicID == nil) == (e.QuizID == nil) {
		return "exercise must reference exactly one of topicId or quizId"
	}
	if e.Prompt == "" {
		return "prompt is required"
	}
	if e.Answer == "" {
		return "answer is required"
	}
	return ""
}

func (s *Server) handleGetExercise(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	exercise, err := s.store.GetExercise(r.Context(), id)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, exercise)
}

func (s *Server) handleCreateExercise(w http.ResponseWriter, r *http.Request) {
	var in model.Exercise
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if msg := validateExercise(in); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	exercise, err := s.store.CreateExercise(r.Context(), in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, exercise)
}

func (s *Server) handleUpdateExercise(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in model.Exercise
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if msg := validateExercise(in); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	exercise, err := s.store.UpdateExercise(r.Context(), id, in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, exercise)
}

func (s *Server) handleDeleteExercise(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.store.DeleteExercise(r.Context(), id); err != nil {
		s.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleAttempt grades a submitted answer, records it and reveals the solution.
func (s *Server) handleAttempt(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var body struct {
		Answer string `json:"answer"`
	}
	if !decodeBody(r, &body) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	attempt, exercise, err := s.store.SubmitAttempt(r.Context(), id, body.Answer)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"attemptId":     attempt.ID,
		"correct":       attempt.IsCorrect,
		"correctAnswer": exercise.Answer,
		"explanation":   exercise.Explanation,
	})
}
