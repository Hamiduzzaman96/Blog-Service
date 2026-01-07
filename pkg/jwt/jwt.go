package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret string
	ttl    time.Duration
}

func New(secret string, ttl time.Duration) *Service {
	return &Service{
		secret: secret,
		ttl:    ttl,
	}
}

func (s *Service) Generate(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{ //payload
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(s.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) /*
				1. JWT header + payload বানাচ্ছে
		        2. Payload হিসেবে claims ব্যবহার করছে
		        3. Sign করার জন্য signing method নির্ধারণ করছে
				HS256 = HMAC SHA-256 (symmetric)Symmetric Key মানে:same secret দিয়ে sign ও verify হয়*/
	return token.SignedString([]byte(s.secret))
}

func (s *Service) Validate(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(s.secret), nil
	})
	/*
	   jwt.Parse(tokenStr, keyFunc) → JWT decode + verify
	   keyFunc → Secret key provide করে backend কে verify করার জন্য
	   t *jwt.Token → parsed token object (header + payload)
	   return []byte(s.secret), nil → symmetric key provide করা HS256 signature check এর জন্য
	*/

	if err != nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return &claims, nil
}
