package response

import (
	"errors"
	"log/slog"
	"net/http"
	"tgiintiz/internal/models"

	"github.com/gin-gonic/gin"
)

type errorBody struct {
	Error string `json:"error"`
}

func JSON(c *gin.Context, status int, body any) {
	c.JSON(status, body)
}

func Error(c *gin.Context, status int, message string, cause error) {
	level := slog.LevelWarn

	if status >= http.StatusInternalServerError {
		level = slog.LevelError
	}

	slog.Log(c.Request.Context(), level, "request failed",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"status", status,
		"error", cause,
	)

	c.AbortWithStatusJSON(status, errorBody{Error: message})
}

func ServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidInput):
		Error(c, http.StatusBadRequest, err.Error(), err)
	case errors.Is(err, models.ErrInvalidCredentials):
		Error(c, http.StatusUnauthorized, models.ErrInvalidCredentials.Error(), err)
	case errors.Is(err, models.ErrEmailTaken):
		Error(c, http.StatusConflict, models.ErrEmailTaken.Error(), err)
	case errors.Is(err, models.ErrNotFound):
		Error(c, http.StatusNotFound, models.ErrNotFound.Error(), err)
	case errors.Is(err, models.ErrInvalidToken):
		Error(c, http.StatusUnauthorized, models.ErrInvalidToken.Error(), err)
	case errors.Is(err, models.ErrForbidden):
		Error(c, http.StatusForbidden, models.ErrForbidden.Error(), err)

	default:
		Error(c, http.StatusInternalServerError, "internal server error", err)
	}
}
