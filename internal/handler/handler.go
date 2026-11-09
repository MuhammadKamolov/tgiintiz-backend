package handler

import (
	"net/http"
	"tgiintiz/internal/handler/middleware"
	v1 "tgiintiz/internal/handler/v1"
	"tgiintiz/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	apiV1 := v1.NewHandler(h.services)

	requireAuth := middleware.Auth(h.services.Auth)
	requireAdmin := func(next http.Handler) http.Handler {
		return requireAuth(middleware.RequireAdmin(next))
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", apiV1.Register)
	mux.HandleFunc("POST /api/v1/auth/login", apiV1.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", apiV1.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", apiV1.Logout)

	mux.Handle("GET /api/v1/users/me", requireAuth(http.HandlerFunc(apiV1.Me)))

	mux.Handle("POST /api/v1/surveys", requireAuth(http.HandlerFunc(apiV1.CreateSurvey)))
	mux.Handle("GET /api/v1/surveys/{id}", requireAuth(http.HandlerFunc(apiV1.GetSurveyByID)))
	mux.Handle("GET /api/v1/surveys", requireAuth(http.HandlerFunc(apiV1.ListMySurveys)))

	mux.Handle("GET /api/v1/admin/surveys", requireAdmin(http.HandlerFunc(apiV1.ListAllSurveys)))
	mux.Handle("PATCH /api/v1/admin/surveys/{id}/status", requireAdmin(http.HandlerFunc(apiV1.UpdateSurveyStatus)))
	mux.Handle("DELETE /api/v1/admin/surveys/{id}", requireAdmin(http.HandlerFunc(apiV1.DeleteSurvey)))

	mux.Handle("POST /api/v1/surveys/{id}/attachments", requireAuth(http.HandlerFunc(apiV1.UploadAttachment)))
	mux.Handle("GET /api/v1/surveys/{id}/attachments", requireAuth(http.HandlerFunc(apiV1.ListAttachments)))
	mux.Handle("DELETE /api/v1/attachments/{id}", requireAuth(http.HandlerFunc(apiV1.DeleteAttachment)))

	fileServer := http.FileServer(http.Dir("./uploads"))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", fileServer))

	return middleware.Logging(middleware.Recover(mux))
}
