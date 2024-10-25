package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/repository/t"
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
		api.Error(w, r, t.Errorf("you cannot modify a different user"), 400)
		return
	}
	auth := r.URL.Query().Get("auth")
	token, err := h.service.LinkFacebook(key, auth, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, token)
}
