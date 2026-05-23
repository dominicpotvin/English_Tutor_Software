package api

import (
	"net/http"

	"englishtutor/internal/model"
)

func (s *Server) handleListVocabulary(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.ListVocabulary(r.Context(),
		optionalID(r, "levelId"),
		optionalID(r, "lessonId"),
		optionalString(r, "category"))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateVocabulary(w http.ResponseWriter, r *http.Request) {
	var in model.Vocabulary
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if in.Term == "" {
		writeError(w, http.StatusBadRequest, "term is required")
		return
	}
	item, err := s.store.CreateVocabulary(r.Context(), in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) handleUpdateVocabulary(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in model.Vocabulary
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	item, err := s.store.UpdateVocabulary(r.Context(), id, in)
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleDeleteVocabulary(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.store.DeleteVocabulary(r.Context(), id); err != nil {
		s.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
