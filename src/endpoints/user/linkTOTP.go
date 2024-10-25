package user

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/repository/t"
)

func (h *Handler) RequestTOTP(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	user, _, err := api.CheckAuth(r)
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	if user != key {
		api.Error(w, r, t.Errorf("you cannot modify a different user"), 400)
		return
	}
	data, err := h.service.RequestTOTP(key, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	api.SuccessJson(w, r, data)
}

func (h *Handler) EnableTOTP(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	user, _, err := api.CheckAuth(r)
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	if user != key {
		api.Error(w, r, t.Errorf("you cannot modify a different user"), 400)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.Error(w, r, t.Errorf("read request body failed: %w", err), 400)
		return
	}
	var totpEnableRequest domain.TotpEnableRequest
	if err := json.Unmarshal(body, &totpEnableRequest); err != nil {
		api.Error(w, r, t.Errorf("decode request body failed: %w", err), 400)
		return
	}
	err = h.service.EnableTOTP(key, totpEnableRequest, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	api.SuccessJson(w, r, nil)
}
