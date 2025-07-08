package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"workout/internal/models"
	"workout/internal/services"
)

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

type UserHandler struct {
	Service *services.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := h.Service.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		models.User
		VerificationCode string `json:"verification_code"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.Service.SignUp(r.Context(), req.User, req.VerificationCode)
	if err != nil {
		if errors.Is(err, models.ErrInvalidVerificationCode) {
			http.Error(w, "Неверный код подтверждения", http.StatusUnauthorized)
			return
		}
		log.Printf("SignUp error: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.Service.SignIn(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UpgradeToTrainer upgrades the authenticated user to trainer role.
func (h *UserHandler) UpgradeToTrainer(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "user id missing", http.StatusUnauthorized)
		return
	}

	if err := h.Service.UpgradeToTrainer(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllClients returns all users with client role.
func (h *UserHandler) GetAllClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.Service.GetAllClients(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

// GetClientsByProgramID lists clients participating in a program.
func (h *UserHandler) GetClientsByProgramID(w http.ResponseWriter, r *http.Request) {
	programID, _ := strconv.Atoi(r.URL.Query().Get(":program_id"))
	if programID == 0 {
		programID, _ = strconv.Atoi(r.URL.Query().Get("program_id"))
	}
	if programID == 0 {
		http.Error(w, "program_id required", http.StatusBadRequest)
		return
	}
	clients, err := h.Service.GetClientsByProgramID(r.Context(), programID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

// DeleteClientFromProgram removes a client's progress from a program.
func (h *UserHandler) DeleteClientFromProgram(w http.ResponseWriter, r *http.Request) {
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
	if err := h.Service.DeleteClientFromProgram(r.Context(), programID, clientID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
