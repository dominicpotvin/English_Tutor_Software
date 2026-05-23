package api

import (
	"net/http"

	"englishtutor/internal/model"
)

// --- Levels ---------------------------------------------------------------

func (s *Server) handleListLevels(w http.ResponseWriter, r *http.Request) {
	levels, err := s.store.ListLevelSummaries(r.Context())
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, levels)
}

func (s *Server) handleGetLevel(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx := r.Context()
	level, err := s.store.GetLevel(ctx, id)
	if err != nil {
		s.fail(w, err)
		return
	}
	lessons, err := s.store.ListLessonSummaries(ctx, id)
	if err != nil {
		s.fail(w, err)
		return
	}
	quizzes, err := s.store.ListQuizSummaries(ctx)
	if err != nil {
		s.fail(w, err)
		return
	}
	levelQuizzes := make([]model.QuizSummary, 0)
	for _, q := range quizzes {
		if q.LevelID != nil && *q.LevelID == id {
			levelQuizzes = append(levelQuizzes, q)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"level":   level,
		"lessons": lessons,
		"quizzes": levelQuizzes,
	})
}

func (s *Server) handleCreateLevel(w http.ResponseWriter, r *http.Request) {
	var in model.Level
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if in.Code == "" || in.Name == "" {
		writeError(w, http.StatusBadRequest, "code and name are required")
		return
	}
	level, err := s.store.CreateLevel(r.Context(), in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, level)
}

func (s *Server) handleUpdateLevel(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in model.Level
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	level, err := s.store.UpdateLevel(r.Context(), id, in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, level)
}

// --- Lessons --------------------------------------------------------------

func (s *Server) handleGetLesson(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx := r.Context()
	lesson, err := s.store.GetLesson(ctx, id)
	if err != nil {
		s.fail(w, err)
		return
	}
	level, err := s.store.GetLevel(ctx, lesson.LevelID)
	if err != nil {
		s.fail(w, err)
		return
	}
	topics, err := s.store.ListTopicSummaries(ctx, id)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"lesson": lesson,
		"level":  level,
		"topics": topics,
	})
}

func (s *Server) handleCreateLesson(w http.ResponseWriter, r *http.Request) {
	var in model.Lesson
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if in.LevelID == 0 || in.Title == "" {
		writeError(w, http.StatusBadRequest, "levelId and title are required")
		return
	}
	lesson, err := s.store.CreateLesson(r.Context(), in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, lesson)
}

func (s *Server) handleUpdateLesson(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in model.Lesson
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	lesson, err := s.store.UpdateLesson(r.Context(), id, in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, lesson)
}

func (s *Server) handleDeleteLesson(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.store.DeleteLesson(r.Context(), id); err != nil {
		s.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Topics ---------------------------------------------------------------

func (s *Server) handleGetTopic(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	topic, err := s.store.GetTopic(r.Context(), id)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, topic)
}

func (s *Server) handleListTopicExercises(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx := r.Context()
	topic, err := s.store.GetTopic(ctx, id)
	if err != nil {
		s.fail(w, err)
		return
	}
	lesson, err := s.store.GetLesson(ctx, topic.LessonID)
	if err != nil {
		s.fail(w, err)
		return
	}
	exercises, err := s.store.ListExercisesByTopic(ctx, id)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"topic":     topic,
		"lesson":    lesson,
		"exercises": exercises,
	})
}

func (s *Server) handleCreateTopic(w http.ResponseWriter, r *http.Request) {
	var in model.Topic
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if in.LessonID == 0 || in.Title == "" {
		writeError(w, http.StatusBadRequest, "lessonId and title are required")
		return
	}
	topic, err := s.store.CreateTopic(r.Context(), in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, topic)
}

func (s *Server) handleUpdateTopic(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in model.Topic
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	topic, err := s.store.UpdateTopic(r.Context(), id, in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, topic)
}

func (s *Server) handleDeleteTopic(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.store.DeleteTopic(r.Context(), id); err != nil {
		s.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
