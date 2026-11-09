package middleware

import (
	"errors"
	"net/http"
	"tgiintiz/internal/handler/response"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFromContext(c)
		if claims == nil {
			response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
			return
		}

		if claims.Role != "admin" {
			response.Error(c, http.StatusForbidden, "forbidden", errors.New("admin role required"))
			return
		}

		c.Next()
	}
}
