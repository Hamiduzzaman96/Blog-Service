package grpc

import (
	"context"

	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"
	"github.com/Hamiduzzaman96/Blog-Service/proto/userpb"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	usecase *usecase.UserUsecase
}

func NewUserhandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	user, err := h.usecase.Register(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		Id:    uint64(user.ID),
		Email: user.Email,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	token, err := h.usecase.Login(req.Email, req.Pasword)
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{
		AccessToken: token,
	}, nil
}
