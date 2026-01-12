package usecase

import (
	"errors"
	"fmt"

	"github.com/Hamiduzzaman96/Blog-Service/internal/domain"
	"github.com/Hamiduzzaman96/Blog-Service/internal/repository"
	"github.com/Hamiduzzaman96/Blog-Service/pkg/jwt"
	"github.com/Hamiduzzaman96/Blog-Service/pkg/redis"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
	jwtSvc   *jwt.Service
	redis    *redis.Client
}

func NewUserUsecase(
	userRepo *repository.UserRepository,
	jwtSvc *jwt.Service,
	redis *redis.Client,
) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
		redis:    redis,
	}
}

func (u *UserUsecase) Register(email, password string) (*domain.User, error) {

	user, err := domain.NewUser(email, password)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hash)

	createdUser, err := u.userRepo.Create(user, user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func (u *UserUsecase) Login(email, password string) (string, error) {
	userModel, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := u.jwtSvc.GenerateAccessToken(userModel.ID, userModel.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	if err := u.redis.SetToken(token, userModel.ID); err != nil {
		return "", fmt.Errorf("failed to store token in redis: %w", err)
	}

	return token, nil
}

func (u *UserUsecase) PromoteToAuthor(userID uint) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.Role == domain.RoleAuthor {
		return errors.New("already an author")
	}

	user.PromoteToAuthor()

	if err := u.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to promote user to author: %w", err)
	}

	return nil
}
