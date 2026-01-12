package usecase

import (
	"errors"
	"time"

	"github.com/Hamiduzzaman96/Blog-Service/internal/domain"
	"github.com/Hamiduzzaman96/Blog-Service/internal/repository"
	"github.com/Hamiduzzaman96/Blog-Service/pkg/rabbitmq"
)

type BlogUsecase struct {
	blogRepo   *repository.BlogRepository
	authorRepo *repository.AuthorRepository
	mq         *rabbitmq.Client
}

func NewBlogUsecase(
	blogRepo *repository.BlogRepository,
	authorRepo *repository.AuthorRepository,
	mq *rabbitmq.Client,
) *BlogUsecase {
	return &BlogUsecase{
		blogRepo:   blogRepo,
		authorRepo: authorRepo,
		mq:         mq,
	}
}

func (b *BlogUsecase) CreatePost(userID uint, title, content string) error {
	author, err := b.authorRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("user is not an author")
	}

	post := &domain.BlogPost{
		AuthorID:  author.ID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if _, err := b.blogRepo.Create(post); err != nil {
		return err
	}

	// async event
	return b.mq.Publish("blog.created", post)
}
