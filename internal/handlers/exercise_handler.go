package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

// UpdateExercise edits an exercise by id.
func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id == 0 {
		id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	}
	if id == 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	var ex models.Exercises
	if err := json.NewDecoder(r.Body).Decode(&ex); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	ex.ID = id

	updated, err := h.Service.UpdateExercise(r.Context(), ex)
	if err != nil {
		if errors.Is(err, models.ErrExerciseNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteExercise removes an exercise by id.
func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id == 0 {
		id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	}
	if id == 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteExercise(r.Context(), id); err != nil {
		if errors.Is(err, models.ErrExerciseNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
