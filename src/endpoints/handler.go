package endpoints

import "pkv/api/src/repository/graph"

type Handler struct {
	db graph.Db
}

func NewHandler(db graph.Db) *Handler {
	return &Handler{db: db}
}
