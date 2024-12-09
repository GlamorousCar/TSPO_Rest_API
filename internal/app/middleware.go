package app

import (
	"net/http"
	"tspo_server/internal/auth"
)

type AuthMiddleware struct {
	jwt *auth.JWTMiddleware
}

func NewAuthMiddleware(jwt *auth.JWTMiddleware) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := m.jwt.ValidateRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if m.jwt.IsBlacklisted(token) {
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
