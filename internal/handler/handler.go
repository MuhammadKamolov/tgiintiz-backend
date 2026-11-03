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

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", apiV1.Register)
	mux.HandleFunc("POST /api/v1/auth/login", apiV1.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", apiV1.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", apiV1.Logout)

	mux.Handle("GET /api/v1/users/me", requireAuth(http.HandlerFunc(apiV1.Me)))

	return middleware.Logging(middleware.Recover(mux))
}
