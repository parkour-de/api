package verband

import (
	verbandService "pkv/api/src/service/verband"
)

type Handler struct {
	service *verbandService.Service
}

func NewHandler(service *verbandService.Service) *Handler {
	return &Handler{service: service}
}
