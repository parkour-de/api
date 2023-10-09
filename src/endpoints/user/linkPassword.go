package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/internal/security"
	"strings"
	"time"
)

// LinkPassword hashes a password and creates a login object for the user and links it to the user

func (h *Handler) LinkPassword(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	key := urlParams.ByName("key")
	password := r.URL.Query().Get("password")
	// check password length
	if len(password) < 8 {
		api.Error(w, r, fmt.Errorf("password too short"), http.StatusBadRequest)
		return
	}
	// check password strength
	if !security.IsStrongPassword(password) {
		api.Error(w, r, fmt.Errorf("password too weak"), http.StatusBadRequest)
		return
	}
	_, err := h.db.Users.Read(key, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read user failed: %w", err), http.StatusBadRequest)
		return
	}
	logins, err := h.db.GetLoginsForUser(key, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read logins failed: %w", err), http.StatusBadRequest)
		return
	}
	for _, login := range logins {
		if login.Provider == "password" {
			api.Error(w, r, fmt.Errorf("password already set"), http.StatusBadRequest)
			return
		}
	}
	nonce := security.MakeNonce()
	login := domain.Login{
		Key:      "",
		Provider: "password",
		Subject:  nonce + ":" + security.HashToken(":password::"+nonce+":"+password),
		Enabled:  true,
		Created:  time.Now(),
	}
	if err = h.db.Logins.Create(&login, r.Context()); err != nil {
		api.Error(w, r, fmt.Errorf("create login failed: %w", err), http.StatusBadRequest)
		return
	}
	if err = h.db.LoginAuthenticatesUser(login, domain.User{Key: key}, r.Context()); err != nil {
		api.Error(w, r, fmt.Errorf("link login to user failed: %w", err), http.StatusBadRequest)
		return
	}
	api.SuccessJson(w, r, nil)
}

func verifyPassword(password, hash string) bool {
	parts := strings.Split(hash, ":")
	if len(parts) != 2 {
		return false
	}
	return security.HashToken(":password::"+parts[0]+":"+password) == parts[1]
}

// VerifyPassword allows the user verify knowledge of the password (for an unknown reason?)
func (h *Handler) VerifyPassword(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	key := urlParams.ByName("key")
	password := r.URL.Query().Get("password")
	_, err := h.db.Users.Read(key, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read user failed: %w", err), http.StatusBadRequest)
		return
	}
	logins, err := h.db.GetLoginsForUser(key, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read logins failed: %w", err), http.StatusBadRequest)
		return
	}
	var login *domain.Login
	for _, l := range logins {
		if l.Provider == "password" {
			login = &l
		}
	}
	login, err = h.db.Logins.Read(login.Key, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read login failed: %w", err), http.StatusBadRequest)
		return
	}
	if !login.Enabled {
		api.Error(w, r, fmt.Errorf("password login not enabled"), http.StatusBadRequest)
		return
	}
	success := verifyPassword(password, login.Subject)
	api.SuccessJson(w, r, success)
}
