package v1

import (
	"errors"
	"net/http"
	"tgiintiz/internal/handler/middleware"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"

	"github.com/gin-gonic/gin"
)

const maxBytes = 10 << 20

func (h *Handler) UploadAttachment(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	requestID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "fail form file", err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to open uploaded file", err)
		return
	}
	defer file.Close()

	input := models.UploadAttachmentInput{
		UserID:    claims.UserID,
		UserRole:  claims.Role,
		RequestID: requestID,
		FileName:  fileHeader.Filename,
		FileSize:  fileHeader.Size,
		Content:   file,
	}

	att, err := h.service.Survey.UploadAttachment(c.Request.Context(), input)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, att)
}

func (h *Handler) ListAttachments(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	requestID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	attachments, err := h.service.Survey.ListAttachments(c.Request.Context(), requestID, claims.UserID, claims.Role)
	if err != nil {
		response.ServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, attachments)
}

func (h *Handler) DeleteAttachment(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	attachmentID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid attachment id", err)
		return
	}

	if err := h.service.Survey.DeleteAttachment(c.Request.Context(), attachmentID, claims.UserID, claims.Role); err != nil {
		response.ServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
