package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"workout/internal/models"
	"workout/internal/services"
)

// ProgramHandler provides HTTP handlers for workout programs.
type ProgramHandler struct {
	Service *services.ProgramService
}

// CreateProgram creates a new workout program.
func (h *ProgramHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	var p models.WorkOutProgram
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if userID, ok := r.Context().Value("user_id").(int); ok {
		p.TrainerID = userID
	} else {
		http.Error(w, "user id missing", http.StatusUnauthorized)
		return
	}

	created, err := h.Service.CreateProgram(r.Context(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// ProgramsByTrainer lists programs for a trainer.
func (h *ProgramHandler) ProgramsByTrainer(w http.ResponseWriter, r *http.Request) {
	trainerIDStr := r.URL.Query().Get("trainer_id")
	trainerID, _ := strconv.Atoi(trainerIDStr)

	if trainerID == 0 {
		if id, ok := r.Context().Value("user_id").(int); ok {
			trainerID = id
		}
	}

	programs, err := h.Service.ProgramsByTrainer(r.Context(), trainerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(programs)
}

// GetProgram returns a workout program by id.
func (h *ProgramHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id == 0 {
		id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	}
	if id == 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	p, err := h.Service.ProgramByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrWorkoutProgramNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// UpdateProgram edits an existing workout program.
func (h *ProgramHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id == 0 {
		id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	}
	if id == 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	var p models.WorkOutProgram
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	p.ID = id

	updated, err := h.Service.UpdateProgram(r.Context(), p)
	if err != nil {
		if errors.Is(err, models.ErrWorkoutProgramNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteProgram removes a program and its days.
func (h *ProgramHandler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id == 0 {
		id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	}
	if id == 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteProgram(r.Context(), id); err != nil {
		if errors.Is(err, models.ErrWorkoutProgramNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
