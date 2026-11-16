package handler

import (
	"net/http"
	"tgiintiz/internal/handler/middleware"
	v1 "tgiintiz/internal/handler/v1"
	"tgiintiz/internal/service"

	"github.com/gin-gonic/gin"
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
	router := gin.New()
	router.Use(middleware.Logging(), middleware.Recover())
	router.Static("/uploads", "./uploads")

	apiV1 := v1.NewHandler(h.services)

	authMiddleware := middleware.Auth(h.services.Auth)
	adminMiddleware := middleware.RequireAdmin()

	api := router.Group("/api/v1")
	{
		// 1. Публичный контур (без авторизации)
		auth := api.Group("/auth")
		{
			auth.POST("/register", apiV1.Register)
			auth.POST("/login", apiV1.Login)
			auth.POST("/refresh", apiV1.Refresh)
			auth.POST("/logout", apiV1.Logout)
		}

		services := api.Group("/services")
		{
			services.GET("", apiV1.ListServices)
			services.GET("/:slug", apiV1.GetServiceBySlug)
		}

		// 2. Клиентский контур (требуется авторизация)
		users := api.Group("/users")
		users.Use(authMiddleware)
		{
			users.GET("/me", apiV1.Me)
		}

		surveys := api.Group("/surveys")
		surveys.Use(authMiddleware)
		{
			surveys.POST("", apiV1.CreateSurvey)
			surveys.GET("", apiV1.ListMySurveys)
			surveys.GET("/:id", apiV1.GetSurveyByID)

			surveys.POST("/:id/services", apiV1.AddServiceToSurvey)
			surveys.GET("/:id/services", apiV1.ListSurveyServices)
			surveys.DELETE("/:id/services", apiV1.RemoveSurveyServices)

			surveys.POST("/:id/attachments", apiV1.UploadAttachment)
			surveys.GET("/:id/attachments", apiV1.ListAttachments)
		}

		attachments := api.Group("/attachments")
		attachments.Use(authMiddleware)
		{
			attachments.DELETE("/:id", apiV1.DeleteAttachment)
		}

		// 3. Админский контур (авторизация + роль admin)
		admin := api.Group("/admin")
		admin.Use(authMiddleware, adminMiddleware)
		{
			adminSurveys := admin.Group("/surveys")
			{
				adminSurveys.GET("", apiV1.ListAllSurveys)
				adminSurveys.PATCH("/:id/status", apiV1.UpdateSurveyStatus)
				adminSurveys.DELETE("/:id", apiV1.DeleteSurvey)
			}

			adminServices := admin.Group("/services")
			{
				adminServices.POST("", apiV1.CreateService)
				adminServices.GET("/:id", apiV1.GetServiceByID)
				adminServices.PATCH("/:id", apiV1.UpdateService)
				adminServices.DELETE("/:id", apiV1.DeleteService)
			}
		}
	}

	return router
}
