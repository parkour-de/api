package authentication

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"strings"
)

type FacebookTokenValidationResponse struct {
	Data struct {
		AppId     string `json:"app_id"`
		ExpiresAt int    `json:"expires_at"`
		IsValid   bool   `json:"is_valid"`
		IssuedAt  int    `json:"issued_at"`
		UserId    string `json:"user_id"`
	} `json:"data"`
}

// Facebook handles the GET /api/facebook endpoint.
func (h *Handler) Facebook(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	auth := r.Header.Get("Authorization")
	format, auth, found := strings.Cut(auth, " ")
	if !found {
		api.Error(w, r, fmt.Errorf("authorization header not correctly formatted"), 400)
		return
	}
	if format != "facebook" {
		api.Error(w, r, fmt.Errorf("authorization header needs to start with 'facebook'"), 400)
		return
	}
	if len(auth) < 1 {
		api.Error(w, r, fmt.Errorf("authorization header contains empty token"), 400)
		return
	}

	tokens, err := h.service.Facebook(auth)
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, tokens)
}
