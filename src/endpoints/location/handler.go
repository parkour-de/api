package location

import (
	"pkv/api/src/domain"
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/photo"
)

type Handler struct {
	db           *graph.Db
	photoService *photo.Service
	em           graph.EntityManager[*domain.Location]
}

func NewHandler(db *graph.Db, photoService *photo.Service, em graph.EntityManager[*domain.Location]) *Handler {
	return &Handler{db: db, photoService: photoService, em: em}
}
