package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Assignment-5/internal/repository"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// GetUsers
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Сбор фильтров
	filters := make(map[string]interface{})
	for key, values := range r.URL.Query() {
		if key == "page" || key == "pageSize" || key == "orderBy" {
			continue
		}
		if len(values) == 0 {
			continue
		}
		if key == "deleted" {
			deletedBool, err := strconv.ParseBool(values[0])
			if err == nil {
				filters["deleted"] = deletedBool
			}
			continue
		}
		filters[key] = values[0]
	}

	orderBy := r.URL.Query().Get("orderBy")

	response, err := h.repo.GetPaginatedUsers(page, pageSize, filters, orderBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCommonFriends
func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id1Str := vars["id1"]
	id2Str := vars["id2"]

	id1, err := strconv.Atoi(id1Str)
	if err != nil {
		http.Error(w, "invalid user ID (id1)", http.StatusBadRequest)
		return
	}
	id2, err := strconv.Atoi(id2Str)
	if err != nil {
		http.Error(w, "invalid user ID (id2)", http.StatusBadRequest)
		return
	}

	friends, err := h.repo.GetCommonFriends(id1, id2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

// SoftDeleteUser
func (h *Handler) SoftDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	result, err := h.repo.DB().Exec("UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "user not found or already deleted", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
