package http

import (
	"encoding/json"
	"net/http"

	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: u,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.Register(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(user) //covert to json and sent to client
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email or password missing", http.StatusBadRequest)
		return
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *UserHandler) PromoteToAuthor(w http.ResponseWriter, r *http.Request) {
	// Get userID from context (set by JWT middleware)
	userIDValue := r.Context().Value("user_id")
	if userIDValue == nil {
		http.Error(w, "user_id missing in context", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		http.Error(w, "invalid user_id type", http.StatusInternalServerError)
		return
	}

	err := h.usecase.PromoteToAuthor(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "user promoted to author successfully",
	})
}
