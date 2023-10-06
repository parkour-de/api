package user

import (
	"pkv/api/src/internal/graph"
)

type Handler struct {
	db *graph.Db
}

func NewHandler(db *graph.Db) *Handler {
	return &Handler{db: db}
}
