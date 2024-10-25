package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/t"
)

type AddToMinecraftWhitelistRequest struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

func (h *Handler) AddUsernameToWhitelist(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	var item AddToMinecraftWhitelistRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, r, t.Errorf("decoding request body failed: %w", err), 400)
		return
	}
	if item.Secret != dpv.ConfigInstance.Auth.MinecraftInviteKey {
		api.Error(w, r, t.Errorf("provided invite key is not correct"), 401)
		return
	}
	err := h.service.AddUsernameToWhitelist(item.Username, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}
