package v1

import "tgiintiz/internal/service"

type Handler struct {
	service *service.Services
}

func NewHandler(service *service.Services) *Handler {
	return &Handler{
		service: service,
	}
}
