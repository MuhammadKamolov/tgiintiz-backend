package v1

import (
	"errors"
	"net/http"
	"tgiintiz/internal/handler/middleware"
	"tgiintiz/internal/handler/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Me(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	user, err := h.service.User.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	resp := newUserResponse(user)

	response.JSON(c, http.StatusOK, resp)

}
