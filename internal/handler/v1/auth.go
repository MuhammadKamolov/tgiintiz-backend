package v1

import (
	"net/http"
	"tgiintiz/internal/handler/response"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	user, err := h.service.Auth.Register(r.Context(), req.toInput())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, newUserResponse(user))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	tokens, err := h.service.Auth.Login(r.Context(), req.toInput())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, newTokenPairResponse(tokens))
}
