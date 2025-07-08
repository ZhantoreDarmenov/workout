package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"workout/internal/models"
	"workout/internal/services"
)

// DayHandler handles workout day requests.
type DayHandler struct {
	Service *services.DayService
}

// DayDetails returns exercises and food for a specific day in a program.
func (h *DayHandler) DayDetails(w http.ResponseWriter, r *http.Request) {
	programID, _ := strconv.Atoi(r.URL.Query().Get(":program_id"))
	if programID == 0 {
		programID, _ = strconv.Atoi(r.URL.Query().Get("program_id"))
	}

	dayNum, _ := strconv.Atoi(r.URL.Query().Get(":day"))
	if dayNum == 0 {
		dayNum, _ = strconv.Atoi(r.URL.Query().Get("day"))
	}

	if programID == 0 || dayNum == 0 {
		http.Error(w, "program_id and day are required", http.StatusBadRequest)
		return
	}

	details, err := h.Service.GetDay(r.Context(), programID, dayNum)
	if err != nil {
		if errors.Is(err, models.ErrDayNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

// CompleteDay marks a day as completed for a client.
func (h *DayHandler) CompleteDay(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID int `json:"client_id"`
		DayID    int `json:"day_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	progress, err := h.Service.CompleteDay(r.Context(), req.ClientID, req.DayID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(progress)
}
func (h *DayHandler) CreateDay(w http.ResponseWriter, r *http.Request) {
	var day models.Days
	if err := json.NewDecoder(r.Body).Decode(&day); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	created, err := h.Service.CreateDay(r.Context(), day)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, models.ErrWorkoutProgramNotFound) || errors.Is(err, models.ErrExerciseNotFound) || errors.Is(err, models.ErrFoodNotFound) {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// DaysByProgram returns all days with details for a program.
func (h *DayHandler) DaysByProgram(w http.ResponseWriter, r *http.Request) {
	programID, _ := strconv.Atoi(r.URL.Query().Get(":program_id"))
	if programID == 0 {
		programID, _ = strconv.Atoi(r.URL.Query().Get("program_id"))
	}
	if programID == 0 {
		http.Error(w, "program_id required", http.StatusBadRequest)
		return
	}

	days, err := h.Service.DaysByProgram(r.Context(), programID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(days)
}

// UpdateDay edits an existing workout day by id.
func (h *DayHandler) UpdateDay(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id == 0 {
		id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	}
	if id == 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	var day models.Days
	if err := json.NewDecoder(r.Body).Decode(&day); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	day.ID = id

	updated, err := h.Service.UpdateDay(r.Context(), day)
	if err != nil {
		if errors.Is(err, models.ErrDayNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrWorkoutProgramNotFound) || errors.Is(err, models.ErrExerciseNotFound) || errors.Is(err, models.ErrFoodNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
