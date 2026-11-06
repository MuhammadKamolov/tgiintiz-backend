package v1

import (
	"net/http"
	"strconv"
	"tgiintiz/internal/handler/middleware"
	"tgiintiz/internal/handler/response"
)

func (h *Handler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var req surveyRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	survey, err := h.service.Survey.Create(r.Context(), claims.UserID, req.toInput())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, newSurveyResponse(survey))
}

func (h *Handler) GetSurveyByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "invalid survey id", err)
		return
	}

	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	survey, err := h.service.Survey.GetByID(r.Context(), id, claims.UserID)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, newSurveyResponse(survey))
}

func (h *Handler) ListMySurveys(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	surveys, err := h.service.Survey.ListByUserID(r.Context(), claims.UserID)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	resp := make([]surveyResponse, 0, len(surveys))

	for _, s := range surveys {
		resp = append(resp, newSurveyResponse(s))
	}

	response.JSON(w, http.StatusOK, resp)
}
