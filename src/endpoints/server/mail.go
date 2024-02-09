package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

type ChangeMailPasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

func (h *Handler) ChangeMailPassword(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	var item ChangeMailPasswordRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&item); err != nil {
		api.Error(w, r, fmt.Errorf("decoding request body failed: %v", err), 400)
		return
	}
	err := h.service.ChangeMailPassword(item.Email, item.OldPassword, item.NewPassword, r.Context())
	if err != nil {
		api.Error(w, r, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}
