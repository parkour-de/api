package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

func (h *Handler) Exists(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	exists, err := h.db.Users.Has(urlParams.ByName("key"), r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("check user exists failed: %w", err), http.StatusBadRequest)
		return
	}
	api.SuccessJson(w, r, exists)
}
