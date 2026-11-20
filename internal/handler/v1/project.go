package v1

import (
	"net/http"
	"strconv"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateProject(c *gin.Context) {
	var req createProjectRequest
	if !decodeJSON(c, &req) {
		return
	}

	project, err := h.service.Project.Create(c.Request.Context(), req.toInput())
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, newProjectResponse(project))
}

func (h *Handler) GetProjectByID(c *gin.Context) {
	projectID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidProjectID, err)
		return
	}

	project, err := h.service.Project.GetByID(c.Request.Context(), projectID)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newProjectResponse(project))
}

func (h *Handler) UpdateProject(c *gin.Context) {
	projectID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidProjectID, err)
		return
	}

	var req updateProjectRequest
	if !decodeJSON(c, &req) {
		return
	}

	project, err := h.service.Project.Update(c.Request.Context(), projectID, req.toInput())
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newProjectResponse(project))
}

func (h *Handler) DeleteProject(c *gin.Context) {
	projectID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidProjectID, err)
		return
	}

	if err := h.service.Project.Delete(c.Request.Context(), projectID); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
func (h *Handler) AddServiceToProject(c *gin.Context) {
	projectID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidProjectID, err)
		return
	}

	var req addProjectServiceRequest
	if !decodeJSON(c, &req) {
		return
	}

	if err := h.service.Project.AddService(c.Request.Context(), projectID, req.ServiceID); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
func (h *Handler) ListProjectServices(c *gin.Context) {
	projectID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidProjectID, err)
		return
	}

	services, err := h.service.Project.ListServicesByProjectID(c.Request.Context(), projectID)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	resp := make([]serviceCardResponse, 0, len(services))
	for _, s := range services {
		resp = append(resp, newServiceCardResponse(s))
	}

	response.JSON(c, http.StatusOK, resp)
}
func (h *Handler) RemoveServiceFromProject(c *gin.Context) {
	projectID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidProjectID, err)
		return
	}

	serviceID, err := strconv.ParseInt(c.Param("serviceId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidServiceID, err)
		return
	}

	if err := h.service.Project.RemoveService(c.Request.Context(), projectID, serviceID); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListProjects(c *gin.Context) {
	lang := c.Query("lang")
	if lang == "" {
		lang = string(models.DefaultLang)
	}

	var (
		projects []*models.Project
		err      error
	)

	if serviceIDStr := c.Query("service_id"); serviceIDStr != "" {
		serviceID, parseErr := strconv.ParseInt(serviceIDStr, 10, 64)
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, msgInvalidServiceID, parseErr)
			return
		}

		projects, err = h.service.Project.ListByServiceID(c.Request.Context(), serviceID, lang)
	} else {
		pagination := parsePagination(c)
		projects, err = h.service.Project.ListActiveByLang(c.Request.Context(), lang, pagination)
	}

	if err != nil {
		response.ServiceError(c, err)
		return
	}

	resp := make([]projectResponse, 0, len(projects))
	for _, project := range projects {
		resp = append(resp, newProjectResponse(project))
	}

	response.JSON(c, http.StatusOK, resp)
}

func (h *Handler) GetProjectBySlug(c *gin.Context) {
	slug := c.Param("slug")

	lang := c.Query("lang")
	if lang == "" {
		lang = string(models.DefaultLang)
	}

	project, err := h.service.Project.GetBySlugAndLang(c.Request.Context(), slug, lang)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, newProjectResponse(project))
}
