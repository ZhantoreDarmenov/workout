package handlers

import (
	"encoding/json"
	"net/http"

	"workout/internal/models"
	"workout/internal/services"
)

// FoodHandler processes HTTP requests for food entries.
type FoodHandler struct {
	Service *services.FoodService
}

func (h *FoodHandler) CreateFood(w http.ResponseWriter, r *http.Request) {
	var f models.Food
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.Service.CreateFood(r.Context(), f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
