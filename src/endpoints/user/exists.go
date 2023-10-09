package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"regexp"
)

func (h *Handler) Exists(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	if err := validateUsername(urlParams.ByName("key")); err != nil {
		api.Error(w, r, fmt.Errorf("invalid username: %w", err), http.StatusBadRequest)
		return
	}
	exists, err := h.db.Users.Has(urlParams.ByName("key"), r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("check user exists failed: %w", err), http.StatusBadRequest)
		return
	}
	api.SuccessJson(w, r, exists)
}

func validateUsername(username string) error {
	/*
		- Valid usernames must match this regex: `^[a-zA-Z0-9_-][a-zA-Z0-9_\-.]{2,29}$`
		  - The first character must be alphanumeric (a-z, A-Z, 0-9), an underscore, or a hyphen.
		  - The remaining characters must be alphanumeric, an underscore, a hyphen, or a period.
		  - The username must be between 3 and 30 characters long.
	*/
	if len(username) < 3 || len(username) > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters long")
	}
	// Regex validation:
	// ^[a-zA-Z0-9_-][a-zA-Z0-9_\-.]{2,29}$
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-][a-zA-Z0-9_\-.]{2,29}$`, username); !matched {
		return fmt.Errorf("username must start with a-z, A-Z, 0-9, _, -, or . but may not start with a period")
	}
	return nil
}
