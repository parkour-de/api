package crud

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/internal/graph"
)

type Handler[T graph.Entity] struct {
	db     graph.DB
	em     graph.EntityManager[T]
	prefix string
}

type IdResponse struct {
	Key string `json:"_key,omitempty" example:"123"`
}

func NewHandler[T graph.Entity](db graph.DB, em graph.EntityManager[T], prefix string) *Handler[T] {
	return &Handler[T]{db, em, prefix}
}

// Create handles the creation of new entities.
func (h *Handler[T]) Create(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	var item T
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, r, fmt.Errorf("decoding request body failed: %w", err), 400)
		return
	}
	err := h.em.Create(item, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("creating entity failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, IdResponse{item.GetKey()})
}

// Read handles the retrieval of entities.
func (h *Handler[T]) Read(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	id := h.prefix + urlParams.ByName("key")
	item, err := h.em.Read(id, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read request failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, item)
}

// Update handles the replacement of existing entities.
func (h *Handler[T]) Update(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	var item T
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, r, fmt.Errorf("decoding request body failed: %w", err), 400)
		return
	}
	err := h.em.Update(item, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("updating entity failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, IdResponse{item.GetKey()})
}

// Delete handles the deletion of entities.
func (h *Handler[T]) Delete(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	id := h.prefix + urlParams.ByName("key")
	var item T
	item.SetKey(id)
	err := h.em.Delete(item, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("deleting entity failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, IdResponse{id})
}
