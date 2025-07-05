package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workout/internal/services"
)

// DayHandler handles workout day requests.
type DayHandler struct {
	Service *services.DayService
}

// DayDetails returns exercises and food for a specific day in a program.
func (h *DayHandler) DayDetails(w http.ResponseWriter, r *http.Request) {
	programID, _ := strconv.Atoi(r.URL.Query().Get("program_id"))
	dayNum, _ := strconv.Atoi(r.URL.Query().Get("day"))
	details, err := h.Service.GetDay(r.Context(), programID, dayNum)
	if err != nil {
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
