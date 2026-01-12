package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hamiduzzaman96/Blog-Service/pkg/jwt"
)

type AuthMiddleware struct {
	jwtService *jwt.Service
}

func NewAuthMiddleware(jwtSvc *jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtSvc}
}

// RequireAuth middleware
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(tokenStr, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtService.Validate(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			http.Error(w, "user_id not found in token", http.StatusUnauthorized)
			return
		}

		// Add user_id to request context
		ctx := context.WithValue(r.Context(), "user_id", uint(userIDFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JWTContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Normally token parse করে user_id বের হবে
		// এখন hardcoded for safety
		var userID uint = 1

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
