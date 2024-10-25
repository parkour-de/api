package authentication

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/repository/t"
	"pkv/api/src/service/user"
)

// Password handles the GET /api/password endpoint.
func (h *Handler) Password(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key, password, ok := r.BasicAuth()
	if !ok {
		api.Error(w, r, t.Errorf("authorization header needs to be RFC 2617 Section 2 compliant"), 400)
		return
	}

	token, err := h.service.Password(key, password, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, token)
}

func (h *Handler) Suggest(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	// Delegate to the service for suggesting a password
	suggestion, err := user.Suggest()
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, suggestion)
}
