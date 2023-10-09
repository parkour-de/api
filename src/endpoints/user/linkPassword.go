package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

func (h *Handler) LinkPassword(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	password := r.URL.Query().Get("password")

	// Delegate to the service for linking the password
	err := h.service.LinkPassword(key, password, r.Context())
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, nil)
}

func (h *Handler) VerifyPassword(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	password := r.URL.Query().Get("password")

	// Delegate to the service for verifying the password
	success, err := h.service.VerifyPassword(key, password, r.Context())
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, success)
}
