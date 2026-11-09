package middleware

import (
	"errors"
	"net/http"
	"strings"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"
	"tgiintiz/internal/service"

	"github.com/gin-gonic/gin"
)

const userCtxKey = "user_claims"

func Auth(auth service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || tokenStr == "" {
			response.Error(c, http.StatusUnauthorized, "missing or invalid bearer token", errors.New("no bearer token in Authorization header"))
			return
		}

		claims, err := auth.Authenticate(c.Request.Context(), tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid or expired token", err)
			return
		}

		c.Set(userCtxKey, claims)

		c.Next()
	}
}

func ClaimsFromContext(c *gin.Context) *models.TokenClaims {
	value, exists := c.Get(userCtxKey)
	if !exists {
		return nil
	}

	claims, ok := value.(*models.TokenClaims)
	if !ok {
		return nil
	}

	return claims
}
