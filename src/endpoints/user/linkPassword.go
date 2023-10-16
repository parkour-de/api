package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

func (h *Handler) LinkPassword(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	user, _, err := api.CheckAuth(r)
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	if user != key {
		api.Error(w, r, fmt.Errorf("you cannot modify a different user"), 400)
		return
	}
	password := r.URL.Query().Get("password")

	if err = h.service.LinkPassword(key, password, r.Context()); err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, nil)
}

func (h *Handler) VerifyPassword(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	user, _, err := api.CheckAuth(r)
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	if user != key {
		api.Error(w, r, fmt.Errorf("you cannot modify a different user"), 400)
		return
	}
	password := r.URL.Query().Get("password")

	success, err := h.service.VerifyPassword(key, password, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, success)
}
