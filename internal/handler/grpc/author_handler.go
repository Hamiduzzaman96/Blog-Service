package grpc

import (
	"context"

	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"
	"github.com/Hamiduzzaman96/Blog-Service/proto/authorpb"
)

type AuthorHandler struct {
	authorpb.UnimplementedAuthorServiceServer
	usecase *usecase.AuthorUsecase
}

func NewAuthorHandler(u *usecase.AuthorUsecase) *AuthorHandler {
	return &AuthorHandler{usecase: u}
}

func (h *AuthorHandler) BecomeAuthor(ctx context.Context, req *authorpb.BecomeAuthorRequest) (*authorpb.AuthorResponse, error) {
	userID := req.UserId

	err := h.usecase.BecomeAuthor(uint(userID))
	if err != nil {
		return nil, err
	}
	return &authorpb.AuthorResponse{
		UserId: userID,
	}, nil
}
