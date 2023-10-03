package user

import (
	"pkv/api/src/domain"
	"pkv/api/src/internal/graph"
)

type Handler struct {
	db graph.DB
	em graph.EntityManager[*domain.User]
}

func NewHandler(db graph.DB, em graph.EntityManager[*domain.User]) *Handler {
	return &Handler{db: db, em: em}
}
