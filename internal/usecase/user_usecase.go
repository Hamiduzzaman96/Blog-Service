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

// Constructor
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

// Register new user
func (u *UserUsecase) Register(email, password string) (*domain.User, error) {
	// Domain-level validation
	user, err := domain.NewUser(email, password)
	if err != nil {
		return nil, err
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hash)

	// Persist
	createdUser, err := u.userRepo.Create(user, user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// Login user and return JWT token
func (u *UserUsecase) Login(email, password string) (string, error) {
	// Fetch user
	userModel, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := u.jwtSvc.GenerateAccessToken(userModel.ID, userModel.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Store token in Redis
	if err := u.redis.SetToken(token, userModel.ID); err != nil {
		return "", fmt.Errorf("failed to store token in redis: %w", err)
	}

	return token, nil
}

// Promote user to Author
func (u *UserUsecase) PromoteToAuthor(userID uint) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.Role == domain.RoleAuthor {
		return errors.New("already an author")
	}

	// Apply domain-level business rule
	user.PromoteToAuthor()

	// Persist update
	return u.userRepo.Update(user)
}
