package middleware

import (
	"context"
	"net/http"
	"strings"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"
	"tgiintiz/internal/service"
)

type contextKey string

const userCtxKey contextKey = "user_payload"

func Auth(auth service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || tokenStr == "" {
				response.Error(w, r, http.StatusUnauthorized, "missing or invalid bearer token", nil)
				return
			}

			payload, err := auth.ParseToken(tokenStr)
			if err != nil {
				response.Error(w, r, http.StatusUnauthorized, "invalid or expired token", err)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserPayload(ctx context.Context) *models.TokenPayload {
	payload, _ := ctx.Value(userCtxKey).(*models.TokenPayload)
	return payload
}
