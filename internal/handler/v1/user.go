package v1

import (
	"net/http"
	"tgiintiz/internal/handler/middleware"
	"tgiintiz/internal/handler/response"
)

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())

	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	user, err := h.service.User.GetByID(r.Context(), claims.UserID)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	resp := newUserResponse(user)

	response.JSON(w, http.StatusOK, resp)

}
