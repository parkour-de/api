package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/internal/security"
	"time"
)

// RequestEmail generates an activation link and sends it via email
// It also creates a login object for the user and links it to the user

func (h *Handler) RequestEmail(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	key := urlParams.ByName("key")
	email := r.URL.Query().Get("email")
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
		if login.Provider == "email" {
			api.Error(w, r, fmt.Errorf("email already requested"), http.StatusBadRequest)
			return
		}
	}
	login := domain.Login{
		Key:      "",
		Provider: "email",
		Subject:  email,
		Enabled:  false,
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
	activationCode := emailActivationCode(login)
	activationLink := fmt.Sprintf("https://parkour-deutschland.de/user/%s/email/%s?code=%s", key, login.Key, activationCode)
	// TODO send email
	log.Println(activationLink)
	api.SuccessJson(w, r, nil)
	// Clients: Send confirmation that an activation link has been sent.
	// Also users should make sure to check if the User id is correct. Clients should show the User id once more.
}

func emailActivationCode(login domain.Login) string {
	token := ":email_activation::" + login.Key + ":" + login.Subject
	activationCode := security.HashToken(token)
	return activationCode
}

// EnableEmail enables the email login for the user
func (h *Handler) EnableEmail(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	_ = urlParams.ByName("key")
	loginId := urlParams.ByName("login")

	// TODO check if authenticated user is the same as the user in the url

	login, err := h.db.Logins.Read(loginId, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read login failed: %w", err), http.StatusBadRequest)
		return
	}
	if login.Provider != "email" {
		api.Error(w, r, fmt.Errorf("invalid provider"), http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("code") != emailActivationCode(*login) {
		api.Error(w, r, fmt.Errorf("invalid activation code"), http.StatusBadRequest)
		return
	}
	if login.Enabled {
		api.Error(w, r, fmt.Errorf("email already enabled"), http.StatusBadRequest)
		return
	}
	login.Enabled = true
	if err := h.db.Logins.Update(login, r.Context()); err != nil {
		api.Error(w, r, fmt.Errorf("update login failed: %w", err), http.StatusBadRequest)
		return
	}
	api.SuccessJson(w, r, nil)
}
