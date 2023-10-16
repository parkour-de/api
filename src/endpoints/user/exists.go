package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

func (h *Handler) Exists(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	username := urlParams.ByName("key")

	// Delegate to the service for checking if the user exists
	exists, err := h.service.Exists(username, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, exists)
}
