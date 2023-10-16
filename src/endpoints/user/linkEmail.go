package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

// RequestEmail generates an activation link and sends it via email
// It also creates a login object for the user and links it to the user

func (h *Handler) RequestEmail(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
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
	email := r.URL.Query().Get("email")
	if err := h.service.RequestEmail(key, email, r.Context()); err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, nil)
}

// EnableEmail enables the email login for the user
func (h *Handler) EnableEmail(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
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
	loginId := urlParams.ByName("login")
	code := r.URL.Query().Get("code")

	// TODO check if authenticated user is the same as the user in the url

	if err := h.service.EnableEmail(loginId, code, r.Context()); err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, nil)
}
