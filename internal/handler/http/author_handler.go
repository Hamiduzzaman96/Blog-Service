package http

import (
	"encoding/json"
	"net/http"

	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"
)

type AuthorHandler struct {
	usecase *usecase.AuthorUsecase
}

func NewAuthorHandler(u *usecase.AuthorUsecase) *AuthorHandler {
	return &AuthorHandler{usecase: u}
}

func (h *AuthorHandler) BecomeAuthor(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uint)
	var req struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := h.usecase.BecomeAuthor(userID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

}
