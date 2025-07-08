package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

// UpdateFood edits a food entry by id.
func (h *FoodHandler) UpdateFood(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id == 0 {
		id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	}
	if id == 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	var f models.Food
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	f.ID = id

	updated, err := h.Service.UpdateFood(r.Context(), f)
	if err != nil {
		if errors.Is(err, models.ErrFoodNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteFood removes a food entry by id.
func (h *FoodHandler) DeleteFood(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id == 0 {
		id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	}
	if id == 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteFood(r.Context(), id); err != nil {
		if errors.Is(err, models.ErrFoodNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
