package handlers

import (
	"encoding/json"
	"net/http"

	"workout/internal/models"
	"workout/internal/services"
)

// ExerciseHandler handles HTTP requests for exercises.
type ExerciseHandler struct {
	Service *services.ExerciseService
}

func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	var ex models.Exercises
	if err := json.NewDecoder(r.Body).Decode(&ex); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.Service.CreateExercise(r.Context(), ex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
