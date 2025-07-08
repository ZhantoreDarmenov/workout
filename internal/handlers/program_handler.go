package handlers

import (
	"encoding/json"
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
