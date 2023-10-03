package query

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

// GetUsers handles the GET /api/users endpoint.
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	// Extract the include parameter from the URL query
	users, err := h.db.GetAllUsers(r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("querying users failed: %w", err), 400)
		return
	}
	api.SuccessJson(w, r, users)
}
