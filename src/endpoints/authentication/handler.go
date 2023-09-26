package authentication

import "pkv/api/src/internal/graph"

type Handler struct {
	db graph.DB
}

func NewHandler(db graph.DB) *Handler {
	return &Handler{db: db}
}
