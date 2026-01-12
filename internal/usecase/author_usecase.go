package usecase

import (
	"errors"

	"github.com/Hamiduzzaman96/Blog-Service/internal/domain"
	"github.com/Hamiduzzaman96/Blog-Service/internal/repository"
)

type AuthorUsecase struct {
	userRepo   *repository.UserRepository
	authorRepo *repository.AuthorRepository
}

func NewAuthorUsecase(
	userRepo *repository.UserRepository,
	authorRepo *repository.AuthorRepository,
) *AuthorUsecase {
	return &AuthorUsecase{
		userRepo:   userRepo,
		authorRepo: authorRepo,
	}
}

func (a *AuthorUsecase) BecomeAuthor(userID uint) error {
	user, err := a.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.Role == domain.RoleAuthor {
		return errors.New("already author")
	}

	user.PromoteToAuthor()

	if err := a.userRepo.Update(user); err != nil {
		return err
	}

	author := domain.NewAuthor(user.ID)
	_, err = a.authorRepo.Create(author)

	return err
}
