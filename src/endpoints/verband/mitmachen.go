package verband

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain/verband"
)

func (h *Handler) Mitmachen(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data verband.MitmachenRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		api.Error(w, r, fmt.Errorf("error decoding request: %w", err), http.StatusBadRequest)
		return
	}

	if err := h.service.Mitmachen(data); err != nil {
		api.Error(w, r, fmt.Errorf("error submitting request: %w", err), http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, map[string]string{"message": "Anfrage erfolgreich abgeschickt!"})
}
