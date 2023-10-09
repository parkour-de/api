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
	if strings.HasPrefix(auth, "facebook ") {
		tokens := strings.SplitN(auth, " ", 2)
		if len(tokens) != 2 {
			api.Error(w, r, fmt.Errorf("authorization header not correctly formatted"), http.StatusBadRequest)
			return
		}
		auth = tokens[1]
		if len(auth) < 1 {
			api.Error(w, r, fmt.Errorf("authorization header contains empty token"), http.StatusBadRequest)
			return
		}
	} else {
		api.Error(w, r, fmt.Errorf("authorization header missing or not correctly prefixed"), http.StatusBadRequest)
		return
	}

	token, err := h.service.Facebook(auth)
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, token)
}
