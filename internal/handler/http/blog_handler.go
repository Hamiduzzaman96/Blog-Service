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
	userID := r.Context().Value("user_id").(uint)
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := h.usecase.CreatePost(userID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
