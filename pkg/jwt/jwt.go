package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewService(secret string, accessExpMin int, refreshExpDay int) *Service {
	return &Service{
		secret:          secret,
		accessTokenTTL:  time.Duration(accessExpMin) * time.Minute,
		refreshTokenTTL: time.Duration(refreshExpDay) * 24 * time.Hour,
	}
}

// Access Token
func (s *Service) GenerateAccessToken(userID uint, role string) (string, error) {
	return s.generateToken(userID, role, s.accessTokenTTL)
}

// Refresh Token
func (s *Service) GenerateRefreshToken(userID uint, role string) (string, error) {
	return s.generateToken(userID, role, s.refreshTokenTTL)
}

func (s *Service) generateToken(userID uint, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// Validate Token
func (s *Service) Validate(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return &claims, nil
}
