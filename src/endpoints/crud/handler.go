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
	ID string `json:"_id,omitempty" example:"item/123"`
}

func NewHandler[T graph.Entity](db graph.DB, em graph.EntityManager[T], prefix string) *Handler[T] {
	return &Handler[T]{db, em, prefix}
}

// Create handles the creation of new entities.
//
//	@Tags		crud
//	@Success	200		{object}	string
//	@Param		request	body		graph.Entity	true	"item to be created"
//	@Success	200		{object}	IdResponse
//	@Router		/api/trainings [post]
//	@Router		/api/locations [post]
//	@Router		/api/users [post]
func (h *Handler[T]) Create(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	var item T
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, fmt.Errorf("decoding request body failed: %w", err), 400)
		return
	}
	err := h.em.Create(item)
	if err != nil {
		api.Error(w, fmt.Errorf("creating entity failed: %w", err), 400)
		return
	}
	jsonMsg, err := json.Marshal(IdResponse{item.GetID()})
	if err != nil {
		api.Error(w, fmt.Errorf("serialising entity failed: %w", err), 400)
		return
	}
	api.Success(w, jsonMsg)
}

// Read handles the retrieval of entities.
//
//	@Tags		crud
//	@Success	200	{object}	string
//	@Param		id	path		string	true	"id of the item to be retrieved"
//	@Success	200	{object}	IdResponse
//	@Router		/api/trainings/{id} [get]
//	@Router		/api/locations/{id} [get]
//	@Router		/api/users/{id} [get]
func (h *Handler[T]) Read(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	id := h.prefix + urlParams.ByName("id")
	item, err := h.em.Read(id)
	if err != nil {
		api.Error(w, fmt.Errorf("read request failed: %w", err), 400)
		return
	}
	jsonMsg, err := json.Marshal(item)
	if err != nil {
		api.Error(w, fmt.Errorf("querying item failed: %w", err), 400)
		return
	}
	api.Success(w, jsonMsg)
}

// Update handles the replacement of existing entities.
//
//	@Tags		crud
//	@Success	200		{object}	string
//	@Param		request	body		graph.Entity	true	"item to be updated"
//	@Success	200		{object}	IdResponse
//	@Router		/api/trainings [put]
//	@Router		/api/locations [put]
//	@Router		/api/users [put]
func (h *Handler[T]) Update(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	var item T
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, fmt.Errorf("decoding request body failed: %w", err), 400)
		return
	}
	err := h.em.Update(item)
	if err != nil {
		api.Error(w, fmt.Errorf("updating entity failed: %w", err), 400)
		return
	}
	jsonMsg, err := json.Marshal(IdResponse{item.GetID()})
	if err != nil {
		api.Error(w, fmt.Errorf("serialising entity failed: %w", err), 400)
		return
	}
	api.Success(w, jsonMsg)
}

// Delete handles the deletion of entities.
//
//	@Tags		crud
//	@Success	200	{object}	string
//	@Param		id	path		string	true	"id of the item to be deleted"
//	@Success	200	{object}	IdResponse
//	@Router		/api/trainings/{id} [delete]
//	@Router		/api/locations/{id} [delete]
//	@Router		/api/users/{id} [delete]
func (h *Handler[T]) Delete(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	id := h.prefix + urlParams.ByName("id")
	var item T
	item.SetID(id)
	err := h.em.Delete(item)
	if err != nil {
		api.Error(w, fmt.Errorf("delete request failed: %w", err), 400)
		return
	}
	jsonMsg, err := json.Marshal(IdResponse{id})
	if err != nil {
		api.Error(w, fmt.Errorf("deleting item failed: %w", err), 400)
		return
	}
	api.Success(w, jsonMsg)
}