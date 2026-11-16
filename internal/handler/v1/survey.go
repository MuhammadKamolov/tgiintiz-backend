package v1

import (
	"errors"
	"net/http"
	"tgiintiz/internal/handler/middleware"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateSurvey(c *gin.Context) {
	var req surveyRequest

	if !decodeJSON(c, &req) {
		return
	}

	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	survey, err := h.service.Survey.Create(c.Request.Context(), claims.UserID, req.toInput())
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, newSurveyResponse(survey))
}

func (h *Handler) GetSurveyByID(c *gin.Context) {
	requestID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	survey, err := h.service.Survey.GetByID(c.Request.Context(), requestID, claims.UserID)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newSurveyResponse(survey))
}

func (h *Handler) ListMySurveys(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	pagination := parsePagination(c)

	surveys, err := h.service.Survey.ListByUserID(c.Request.Context(), claims.UserID, pagination)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	resp := make([]surveyResponse, 0, len(surveys))
	for _, survey := range surveys {
		resp = append(resp, newSurveyResponse(survey))
	}

	response.JSON(c, http.StatusOK, resp)
}

func (h *Handler) ListAllSurveys(c *gin.Context) {
	var (
		surveys []*models.SurveyRequest
		err     error
	)

	status := c.Query("status")
	pagination := parsePagination(c)

	if status != "" {
		surveys, err = h.service.Survey.ListByStatus(c.Request.Context(), status, pagination)
	} else {
		surveys, err = h.service.Survey.ListAll(c.Request.Context(), pagination)
	}
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	resp := make([]surveyResponse, 0, len(surveys))
	for _, survey := range surveys {
		resp = append(resp, newSurveyResponse(survey))
	}

	response.JSON(c, http.StatusOK, resp)
}

func (h *Handler) UpdateSurveyStatus(c *gin.Context) {
	requestID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	var req updateSurveyStatusRequest

	if !decodeJSON(c, &req) {
		return
	}

	survey, err := h.service.Survey.UpdateStatus(c.Request.Context(), requestID, req.Status)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newSurveyResponse(survey))
}

func (h *Handler) DeleteSurvey(c *gin.Context) {
	requestID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	if err := h.service.Survey.Delete(c.Request.Context(), requestID); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) AddServiceToSurvey(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	requestID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidRequestID, err)
		return
	}

	var req addSurveyServiceRequest
	if !decodeJSON(c, &req) {
		return
	}

	if err := h.service.Survey.AddService(c.Request.Context(), requestID, req.ServiceID, *claims); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) ListSurveyServices(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	requestID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidRequestID, err)
		return
	}

	serviceCards, err := h.service.Survey.ListServices(c.Request.Context(), requestID, *claims)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	resp := make([]serviceCardResponse, 0, len(serviceCards))
	for _, card := range serviceCards {
		resp = append(resp, newServiceCardResponse(card))
	}

	response.JSON(c, http.StatusOK, resp)

}

func (h *Handler) RemoveSurveyServices(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	requestID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidRequestID, err)
		return
	}

	if err := h.service.Survey.RemoveServices(c.Request.Context(), requestID, *claims); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
