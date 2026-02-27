package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   *int   `json:"age,omitempty"`
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   *int   `json:"age,omitempty"`
}

// GetAllUsers
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserUsecase.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, users)
}

// GetUserByID
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.UserUsecase.GetUserByID(r.Context(), id)
	if err != nil {
		if err.Error() == "user with id "+idStr+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// CreateUser
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.UserUsecase.CreateUser(r.Context(), req.Name, req.Email, req.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, http.StatusCreated, map[string]int{"id": id})
}

// UpdateUser
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.UserUsecase.UpdateUser(r.Context(), id, req.Name, req.Email, req.Age)
	if err != nil {
		if err.Error() == "user with id "+idStr+" does not exist" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.UserUsecase.DeleteUser(r.Context(), id)
	if err != nil {
		if err.Error() == "user with id "+idStr+" does not exist" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Healthcheck
func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// respondJSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
