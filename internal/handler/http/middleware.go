package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hamiduzzaman96/Blog-Service/pkg/jwt"
)

type AuthMiddleware struct {
	jwtSvc *jwt.Service
}

func NewAuthMiddleware(jwtSvc *jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{jwtSvc: jwtSvc}
}

func (a *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer")
		claims, err := a.jwtSvc.Validate(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", uint((*claims)["user_id"].(float64)))
		ctx = context.WithValue(ctx, "role", (*claims)["role"].(string))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rRole := r.Context().Value("role").(string)
		if rRole != role {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
