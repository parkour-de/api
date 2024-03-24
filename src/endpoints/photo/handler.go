package photo

import (
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/photo"
)

type Handler struct {
	service *photo.Service
}

type PhotoEntityHandler[T graph.PhotoEntity] struct {
	service *photo.Service
	em      graph.EntityManager[T]
}

func NewHandler(service *photo.Service) *Handler {
	return &Handler{service: service}
}

func NewPhotoEntityHandler[T graph.PhotoEntity](service *photo.Service, em graph.EntityManager[T]) *PhotoEntityHandler[T] {
	return &PhotoEntityHandler[T]{service: service, em: em}
}
