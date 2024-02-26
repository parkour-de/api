package photo

import (
	"pkv/api/src/service/photo"
)

type Handler struct {
	service *photo.Service
}

func NewHandler(service *photo.Service) *Handler {
	return &Handler{service: service}
}
