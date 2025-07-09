package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"workout/internal/models"
	"workout/internal/services"
)

// InviteHandler exposes HTTP endpoints for program invites.
type InviteHandler struct {
	Service *services.InviteService
}

func (h *InviteHandler) InviteClient(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProgramID  int    `json:"work_out_program_id"`
		Email      string `json:"email"`
		Message    string `json:"message"`
		AccessDays int    `json:"access_days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	inv, err := h.Service.InviteClient(r.Context(), req.ProgramID, req.Email, req.Message, req.AccessDays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inv)
}

func (h *InviteHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	clientID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "user id missing", http.StatusUnauthorized)
		return
	}
	inv, err := h.Service.AcceptInvite(r.Context(), req.Token, clientID)
	if err != nil {
		if err == models.ErrInviteNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func (h *InviteHandler) UpdateAccess(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessDays int `json:"access_days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	programID, _ := strconv.Atoi(r.URL.Query().Get(":program_id"))
	if programID == 0 {
		programID, _ = strconv.Atoi(r.URL.Query().Get("program_id"))
	}
	clientID, _ := strconv.Atoi(r.URL.Query().Get(":client_id"))
	if clientID == 0 {
		clientID, _ = strconv.Atoi(r.URL.Query().Get("client_id"))
	}
	if programID == 0 || clientID == 0 {
		http.Error(w, "program_id and client_id required", http.StatusBadRequest)
		return
	}
	inv, err := h.Service.UpdateAccess(r.Context(), programID, clientID, req.AccessDays)
	if err != nil {
		if err == models.ErrInviteNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}


func (h *InviteHandler) ProgramFromInvite(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}
	program, err := h.Service.GetProgramFromInvite(r.Context(), token)
	if err != nil {
		if err == models.ErrInviteNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(program)
}

