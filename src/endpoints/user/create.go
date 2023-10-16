package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/service/user"
	"time"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	username := r.URL.Query().Get("username")
	userType := r.URL.Query().Get("type")

	key, err := h.service.Create(key, username, userType, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	token := user.HashedUserToken("a", key, time.Now().Add(time.Minute*30).Unix())

	api.SuccessJson(w, r, token)
}

func (h *Handler) Claim(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")

	err := h.service.Claim(key, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, nil)
}
