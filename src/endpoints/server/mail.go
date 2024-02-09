package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
		http.Error(w, fmt.Sprintf("decoding request body failed: %v", err), http.StatusBadRequest)
		return
	}
	err := h.service.ChangeMailPassword(item.Email, item.OldPassword, item.NewPassword, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
