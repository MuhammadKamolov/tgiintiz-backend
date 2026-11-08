package v1

import (
	"errors"
	"net/http"
	"tgiintiz/internal/handler/middleware"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"
)

const msgInvalidSurveyID = "invalid survey id"

func (h *Handler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var req surveyRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
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
	id, err := parseID(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
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
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	pagination := parsePagination(r)

	surveys, err := h.service.Survey.ListByUserID(r.Context(), claims.UserID, pagination)
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

func (h *Handler) ListAllSurveys(w http.ResponseWriter, r *http.Request) {
	var (
		survey []*models.SurveyRequest
		err    error
	)

	status := r.URL.Query().Get("status")
	pagination := parsePagination(r)

	if status != "" {
		survey, err = h.service.Survey.ListByStatus(r.Context(), status, pagination)
	} else {
		survey, err = h.service.Survey.ListAll(r.Context(), pagination)
	}
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	resp := make([]surveyResponse, 0, len(survey))

	for _, s := range survey {
		resp = append(resp, newSurveyResponse(s))
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateSurveyStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	var req updateSurveyStatusRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	survey, err := h.service.Survey.UpdateStatus(r.Context(), id, req.Status)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, newSurveyResponse(survey))
}

func (h *Handler) DeleteSurvey(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	if err := h.service.Survey.Delete(r.Context(), id); err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
