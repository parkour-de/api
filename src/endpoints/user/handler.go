package user

import (
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/user"
)

type Handler struct {
	db      *graph.Db
	service *user.Service
}

func NewHandler(db *graph.Db, service *user.Service) *Handler {
	return &Handler{db: db, service: service}
}
