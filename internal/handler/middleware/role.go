package middleware

import (
	"errors"
	"net/http"
	"slices"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"

	"github.com/gin-gonic/gin"
)

func RequireRoles(allowed ...models.RoleID) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFromContext(c)
		if claims == nil {
			response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
			return
		}

		if !slices.Contains(allowed, claims.RoleID) {
			response.Error(c, http.StatusForbidden, "forbidden", errors.New("insufficient permissions"))
			return
		}

		c.Next()
	}
}
