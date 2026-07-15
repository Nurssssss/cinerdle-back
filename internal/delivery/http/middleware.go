package http

import (
	"cinerdle-back/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwtManager jwt.JwtManager
}

func NewAuthMiddleware(jwtManager jwt.JwtManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (aum *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		userID, err := aum.jwtManager.Verify(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}
