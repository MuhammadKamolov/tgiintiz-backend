package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"
	"tgiintiz/internal/service"
)

type contextKey string

const userCtxKey contextKey = "user_claims"

func Auth(auth service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || tokenStr == "" {
				response.Error(w, r, http.StatusUnauthorized, "missing or invalid bearer token", errors.New("no bearer token in Authorization header"))
				return
			}

			claims, err := auth.Authenticate(r.Context(), tokenStr)
			if err != nil {
				response.Error(w, r, http.StatusUnauthorized, "invalid or expired token", err)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) *models.TokenClaims {
	claims, _ := ctx.Value(userCtxKey).(*models.TokenClaims)
	return claims
}
