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
	json.NewDecoder(r.Body).Decode(&req)

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
	json.NewDecoder(r.Body).Decode(&req)

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
