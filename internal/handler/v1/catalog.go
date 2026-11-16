package v1

import (
	"net/http"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateService(c *gin.Context) {
	var req createServiceRequest

	if !decodeJSON(c, &req) {
		return
	}

	service, err := h.service.Catalog.Create(c.Request.Context(), req.toInput())
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, newServiceResponse(service))
}

func (h *Handler) GetServiceByID(c *gin.Context) {
	serviceID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidServiceID, err)
		return
	}

	service, err := h.service.Catalog.GetByID(c.Request.Context(), serviceID)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newServiceResponse(service))
}

func (h *Handler) UpdateService(c *gin.Context) {
	serviceID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidServiceID, err)
		return
	}

	var req updateServiceRequest
	if !decodeJSON(c, &req) {
		return
	}

	service, err := h.service.Catalog.Update(c.Request.Context(), serviceID, req.toInput())
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newServiceResponse(service))
}

func (h *Handler) DeleteService(c *gin.Context) {
	serviceID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidServiceID, err)
		return
	}

	if err := h.service.Catalog.Delete(c.Request.Context(), serviceID); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListServices(c *gin.Context) {
	lang := c.Query("lang")
	if lang == "" {
		lang = string(models.DefaultLang)
	}

	serviceCards, err := h.service.Catalog.ListByLang(c.Request.Context(), lang)
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

func (h *Handler) GetServiceBySlug(c *gin.Context) {
	slug := c.Param("slug")
	lang := c.Query("lang")
	if lang == "" {
		lang = string(models.DefaultLang)
	}

	service, err := h.service.Catalog.GetBySlugAndLang(c.Request.Context(), slug, lang)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newServiceResponse(service))
}
