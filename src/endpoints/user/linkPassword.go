package user

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/repository/t"
)

func (h *Handler) Password(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.loginWithPassword(w, r, key)
		return
	}

	h.linkPassword(w, r, key)
}

func (h *Handler) linkPassword(w http.ResponseWriter, r *http.Request, key string) {
	user, _, err := api.CheckAuth(r)
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	if user != key {
		api.Error(w, r, t.Errorf("you cannot modify a different user"), 400)
		return
	}

	password, err := extractPassword(r)
	if err != nil {
		api.Error(w, r, t.Errorf("invalid request body: %w", err), 400)
		return
	}

	token, err := h.service.LinkPassword(key, password, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, token)
}

func (h *Handler) loginWithPassword(w http.ResponseWriter, r *http.Request, key string) {
	password, err := extractPassword(r)
	if err != nil {
		api.Error(w, r, t.Errorf("invalid request body: %w", err), 400)
		return
	}

	token, err := h.service.Password(key, password, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, token)
	return
}

func extractPassword(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", t.Errorf("request body missing")
	}
	var requestBody struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return "", err
	}

	return requestBody.Password, nil
}

func (h *Handler) VerifyPassword(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	key := urlParams.ByName("key")
	user, _, err := api.CheckAuth(r)
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}
	if user != key {
		api.Error(w, r, t.Errorf("you cannot modify a different user"), 400)
		return
	}
	password := r.URL.Query().Get("password")

	success, err := h.service.VerifyPassword(key, password, r.Context())
	if err != nil {
		api.Error(w, r, err, 400)
		return
	}

	api.SuccessJson(w, r, success)
}
