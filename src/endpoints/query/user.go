package query

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
)

// GetUsers handles the GET /api/users endpoint.
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	query := r.URL.Query()
	skip, err := api.ParseInt(query.Get("skip"))
	if err != nil {
		api.Error(w, r, fmt.Errorf("invalid skip: %w", err), 400)
		return
	}
	limit, err := api.ParseInt(query.Get("limit"))
	if err != nil {
		api.Error(w, r, fmt.Errorf("invalid limit: %w", err), 400)
		return
	}
	queryOptions := domain.UserQueryOptions{
		Key:      query.Get("key"),
		Name:     query.Get("name"),
		Type:     query.Get("type"),
		Text:     query.Get("text"),
		Language: query.Get("language"),
		Include:  api.MakeSet(query.Get("include")),
		Skip:     skip,
		Limit:    limit,
	}
	users, err := h.db.GetFilteredUsers(queryOptions, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("querying users failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, users)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	item, err := h.db.Users.Read(key, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read request failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, item)
}
