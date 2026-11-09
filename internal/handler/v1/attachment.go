package v1

import (
	"errors"
	"net/http"
	"tgiintiz/internal/handler/middleware"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"
)

const maxBytes = 10 << 20

func (h *Handler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	requestID, err := parseID(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "fail form file", err)
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

	att, err := h.service.Survey.UploadAttachment(r.Context(), input)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, att)
}

func (h *Handler) ListAttachments(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	requestID, err := parseID(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, msgInvalidSurveyID, err)
		return
	}

	attachments, err := h.service.Survey.ListAttachments(r.Context(), requestID, claims.UserID, claims.Role)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, attachments)
}

func (h *Handler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, r, http.StatusUnauthorized, "unauthorized", errors.New("claims not found in context"))
		return
	}

	attachmentID, err := parseID(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "invalid attachment id", err)
		return
	}

	if err := h.service.Survey.DeleteAttachment(r.Context(), attachmentID, claims.UserID, claims.Role); err != nil {
		response.ServiceError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
