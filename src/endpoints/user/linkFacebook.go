package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

// LinkFacebook attaches a facebook subject (sub) to a user
func (h *Handler) LinkFacebook(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	auth := r.URL.Query().Get("auth")
	err := h.service.LinkFacebook(key, auth, r.Context())
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, nil)
}
