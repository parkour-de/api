package query

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

// GetPages handles the GET /api/pages endpoint.
//
//	@Summary		Get a list of pages
//	@Description	Returns a list of pages.
//	@Tags			users
//	@Success		200	{array}	domain.Page
//	@Router			/api/pages [get]
func (h *Handler) GetPages(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	// Extract the include parameter from the URL query
	pages, err := h.db.GetAllPages()
	if err != nil {
		api.Error(w, fmt.Errorf("querying users failed: %w", err), 400)
		return
	}
	jsonMsg, err := json.Marshal(pages)
	if err != nil {
		api.Error(w, fmt.Errorf("querying users failed: %w", err), 400)
		return
	}

	api.Success(w, jsonMsg)
}
