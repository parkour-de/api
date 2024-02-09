package server

import "pkv/api/src/service/server"

type Handler struct {
	service *server.Service
}

func NewHandler(service *server.Service) *Handler {
	return &Handler{service: service}
}
