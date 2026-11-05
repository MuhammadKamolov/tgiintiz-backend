package v1

import (
	"net/http"
	"tgiintiz/internal/handler/middleware"
	"tgiintiz/internal/handler/response"
)

func (h *Handler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var req surveyRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	claims := middleware.ClaimsFromContext(r.Context())

	survey, err := h.service.Survey.Create(r.Context(), claims.UserID, req.toInput())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, newSurveyResponse(survey))
}
