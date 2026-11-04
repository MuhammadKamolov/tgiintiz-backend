package middleware

import (
	"net/http"
	"tgiintiz/internal/handler/response"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := ClaimsFromContext(r.Context())
		if claims == nil {
			response.Error(w, r, http.StatusUnauthorized, "unauthorized", nil)
			return
		}

		if claims.Role != "admin" {
			response.Error(w, r, http.StatusForbidden, "forbidden: admin access required", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
