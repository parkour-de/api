package query

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

// GetPages handles the GET /api/pages endpoint.
func (h *Handler) GetPages(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	// Extract the include parameter from the URL query
	pages, err := h.db.GetAllPages(r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("querying pages failed: %w", err), 400)
		return
	}

	api.SuccessJson(w, r, pages)
}
