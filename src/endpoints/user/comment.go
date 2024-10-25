package user

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/repository/t"
)

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	var item domain.Comment
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, r, t.Errorf("decoding request body failed: %w", err), 400)
		return
	}
	var err error
	item.Author, _, err = api.RequireUserAdmin(item.Author, r, h.db)
	if err != nil {
		api.Error(w, r, t.Errorf("cannot comment as %s: %w", item.Author, err), 400)
		return
	}
	key := urlParams.ByName("key")
	if err = h.service.AddComment(key, item.Author, item.Title, item.Text, r.Context()); err != nil {
		api.Error(w, r, err, 400)
		return
	}
	api.SuccessJson(w, r, nil)
}

func (h *Handler) EditComment(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	title := r.URL.Query().Get("title")
	var item domain.Comment
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, r, t.Errorf("decoding request body failed: %w", err), 400)
		return
	}
	var err error
	item.Author, _, err = api.RequireUserAdmin(item.Author, r, h.db)
	if err != nil {
		api.Error(w, r, t.Errorf("cannot edit comments of %s: %w", item.Author, err), 400)
		return
	}
	key := urlParams.ByName("key")
	if err = h.service.EditComment(key, item.Author, title, item.Title, item.Text, r.Context()); err != nil {
		api.Error(w, r, err, 400)
		return
	}
	api.SuccessJson(w, r, nil)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")
	var err error
	author, _, err = api.RequireUserAdmin(author, r, h.db)
	if err != nil {
		api.Error(w, r, t.Errorf("cannot delete comment of %s: %w", author, err), 400)
		return
	}
	key := urlParams.ByName("key")
	if err = h.service.DeleteComment(key, author, title, r.Context()); err != nil {
		api.Error(w, r, err, 400)
		return
	}
	api.SuccessJson(w, r, nil)
}
