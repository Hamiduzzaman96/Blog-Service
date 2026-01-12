package http

import (
	"encoding/json"
	"net/http"

	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"
)

type BlogHandler struct {
	usecase *usecase.BlogUsecase
}

func NewBolgHandler(u *usecase.BlogUsecase) *BlogHandler {
	return &BlogHandler{usecase: u}
}

func (h *BlogHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user_id")
	roleVal := r.Context().Value("user_role")

	if userIDVal == nil || roleVal == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := userIDVal.(uint)
	role := roleVal.(string)

	if role != "AUTHOR" {
		http.Error(w, "user is not an author", http.StatusForbidden)
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.usecase.CreatePost(userID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
