package handlers

import (
	"encoding/json"
	"net/http"

	"workout/internal/services"
)

// AnalyticsHandler exposes endpoints for analytics operations.
type AnalyticsHandler struct {
	Service *services.AnalyticsService
}

func (h *AnalyticsHandler) TrainerAnalytics(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "user id missing", http.StatusUnauthorized)
		return
	}
	data, err := h.Service.TrainerAnalytics(r.Context(), trainerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
