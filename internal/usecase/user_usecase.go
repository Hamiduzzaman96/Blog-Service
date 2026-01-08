package usecase

import (
	"errors"

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

// Register

func (u *UserUsecase) Register(email, password string) (*domain.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	user, err := domain.NewUser(email, password)
	if err != nil {
		return nil, err
	}

	user.Password = string(hash)
	return u.userRepo.Create(user, string(hash))
}

// Login

func (u *UserUsecase) Login(email, password string) (string, error) {
	userModel, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	// password verification

	if err := bcrypt.CompareHashAndPassword(
		[]byte(userModel.Password),
		[]byte(password),
	); err != nil {
		return "", errors.New("Invalid creditials")
	}

	token, err := u.jwtSvc.Generate(userModel.ID, userModel.Role)
	if err != nil {
		return "", err
	}

	if err := u.redis.SetToken(token, userModel.ID); err != nil {
		return "", err
	}

	return token, nil
}
