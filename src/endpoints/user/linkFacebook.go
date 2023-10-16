package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

// LinkFacebook attaches a facebook subject (sub) to a user
func (h *Handler) LinkFacebook(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
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
	auth := r.URL.Query().Get("auth")
	if err = h.service.LinkFacebook(key, auth, r.Context()); err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, nil)
}
