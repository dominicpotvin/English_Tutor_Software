package api

import "net/http"

func (s *Server) handleProgress(w http.ResponseWriter, r *http.Request) {
	levels, err := s.store.LevelProgress(r.Context())
	if err != nil {
		s.fail(w, err)
		return
	}
	totals := map[string]int{
		"totalExercises":    0,
		"masteredExercises": 0,
		"totalAttempts":     0,
		"correctAttempts":   0,
	}
	for _, lp := range levels {
		totals["totalExercises"] += lp.TotalExercises
		totals["masteredExercises"] += lp.MasteredExercises
		totals["totalAttempts"] += lp.TotalAttempts
		totals["correctAttempts"] += lp.CorrectAttempts
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"levels": levels,
		"totals": totals,
	})
}
