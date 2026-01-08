package grpc

import (
	"context"

	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"
	"github.com/Hamiduzzaman96/Blog-Service/proto/blogpb"
)

type Bloghandler struct {
	blogpb.UnimplementedBlogServiceServer
	usecase *usecase.BlogUsecase
}

func NewBlogHandler(u *usecase.BlogUsecase) *Bloghandler {
	return &Bloghandler{usecase: u}
}

func (h *Bloghandler) CreatePost(ctx context.Context, req *blogpb.CreatePostRequest) (*blogpb.BlogResponse, error) {
	err := h.usecase.CreatePost(uint(req.AuthorId), req.Title, req.Content)
	if err != nil {
		return nil, err
	}
	return &blogpb.BlogResponse{
		Title:   req.Title,
		Content: req.Content,
	}, nil
}
