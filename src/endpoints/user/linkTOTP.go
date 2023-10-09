package user

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
)

func (h *Handler) RequestTOTP(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	data, err := h.service.RequestTOTP(key, r.Context())
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}
	api.SuccessJson(w, r, data)
}

func (h *Handler) EnableTOTP(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.Error(w, r, fmt.Errorf("read request body failed: %w", err), http.StatusBadRequest)
		return
	}
	var totpEnableRequest domain.TotpEnableRequest
	if err := json.Unmarshal(body, &totpEnableRequest); err != nil {
		api.Error(w, r, fmt.Errorf("decode request body failed: %w", err), http.StatusBadRequest)
		return
	}
	err = h.service.EnableTOTP(key, totpEnableRequest, r.Context())
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}
	api.SuccessJson(w, r, nil)
}
