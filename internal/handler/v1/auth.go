package v1

import (
	"net/http"
	"tgiintiz/internal/handler/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest

	if !decodeJSON(c, &req) {
		return
	}

	user, err := h.service.Auth.Register(c.Request.Context(), req.toInput())
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, newUserResponse(user))
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest

	if !decodeJSON(c, &req) {
		return
	}

	tokens, err := h.service.Auth.Login(c.Request.Context(), req.toInput())
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newTokenPairResponse(tokens))
}

func (h *Handler) Refresh(c *gin.Context) {
	var req refreshTokenRequest

	if !decodeJSON(c, &req) {
		return
	}

	tokens, err := h.service.Auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newTokenPairResponse(tokens))
}

func (h *Handler) Logout(c *gin.Context) {
	var req refreshTokenRequest

	if !decodeJSON(c, &req) {
		return
	}

	if err := h.service.Auth.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
