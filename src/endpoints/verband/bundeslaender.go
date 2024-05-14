package verband

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

func (h *Handler) GetBundeslaender(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	list, err := h.service.BundeslandInfo(r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	api.SuccessJson(w, r, list)
}
