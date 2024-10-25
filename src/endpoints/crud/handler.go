package crud

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/repository/graph"
	"pkv/api/src/repository/t"
)

type Handler[T graph.Entity] struct {
	db *graph.Db
	em graph.EntityManager[T]
}

type KeyResponse struct {
	Key string `json:"_key,omitempty" example:"123"`
}

func NewHandler[T graph.Entity](db *graph.Db, em graph.EntityManager[T]) *Handler[T] {
	return &Handler[T]{db, em}
}

// Create handles the creation of new entities.
func (h *Handler[T]) Create(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	_, err := api.RequireGlobalAdmin(r, h.db)
	if err != nil {
		api.Error(w, r, t.Errorf("cannot perform CREATE operation: %w", err), 400)
		return
	}
	var item T
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, r, t.Errorf("decoding request body failed: %w", err), 400)
		return
	}
	err = h.em.Create(item, r.Context())
	if err != nil {
		api.Error(w, r, t.Errorf("creating entity failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, KeyResponse{item.GetKey()})
}

// Read handles the retrieval of entities.
func (h *Handler[T]) Read(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	_, err := api.RequireGlobalAdmin(r, h.db)
	if err != nil {
		api.Error(w, r, t.Errorf("cannot perform READ operation: %w", err), 400)
		return
	}
	key := urlParams.ByName("key")
	item, err := h.em.Read(key, r.Context())
	if err != nil {
		api.Error(w, r, t.Errorf("read request failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, item)
}

// Update handles the replacement of existing entities.
func (h *Handler[T]) Update(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	_, err := api.RequireGlobalAdmin(r, h.db)
	if err != nil {
		api.Error(w, r, t.Errorf("cannot perform UPDATE operation: %w", err), 400)
		return
	}
	item, err := h.PostBody(r)
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	err = h.em.Update(item, r.Context())
	if err != nil {
		api.Error(w, r, t.Errorf("updating entity failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, KeyResponse{item.GetKey()})
}

func (h *Handler[T]) PostBody(r *http.Request) (T, error) {
	var item T
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		return item, t.Errorf("decoding request body failed: %w", err)
	}
	return item, nil
}

// Delete handles the deletion of entities.
func (h *Handler[T]) Delete(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	_, err := api.RequireGlobalAdmin(r, h.db)
	if err != nil {
		api.Error(w, r, t.Errorf("cannot perform DELETE operation: %w", err), 400)
		return
	}
	key := urlParams.ByName("key")
	var item T
	item.SetKey(key)
	err = h.em.Delete(item, r.Context())
	if err != nil {
		api.Error(w, r, t.Errorf("deleting entity failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, KeyResponse{key})
}
