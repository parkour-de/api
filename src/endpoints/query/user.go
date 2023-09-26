package query

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
)

// GetUsers handles the GET /api/users endpoint.
//
//	@Summary		Get a list of users
//	@Description	Returns a list of users.
//	@Tags			users
//	@Success		200	{array}	domain.User
//	@Router			/api/users [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	// Extract the include parameter from the URL query
	users, err := h.db.GetAllUsers()
	if err != nil {
		api.Error(w, fmt.Errorf("querying users failed: %w", err), 400)
		return
	}
	jsonMsg, err := json.Marshal(users)
	if err != nil {
		api.Error(w, fmt.Errorf("querying users failed: %w", err), 400)
		return
	}

	api.Success(w, jsonMsg)
}
