package user

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pquerna/otp/totp"
	"image/png"
	"io"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"time"
)

// RequestTOTP generates a totp key and returns it as a base64 encoded png image
// It also creates a login object for the user and links it to the user

func (h *Handler) RequestTOTP(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	key := urlParams.ByName("key")
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
		if login.Provider == "totp" {
			api.Error(w, r, fmt.Errorf("totp already requested"), http.StatusBadRequest)
			return
		}
	}
	otp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "parkour-deutschland.de",
		AccountName: key, // TODO check if authenticated user is the same as the user in the url
	})
	if err != nil {
		api.Error(w, r, fmt.Errorf("generate totp key failed: %w", err), http.StatusBadRequest)
		return
	}
	var buf bytes.Buffer
	img, err := otp.Image(200, 200)
	if err != nil {
		api.Error(w, r, fmt.Errorf("generate totp image failed: %w", err), http.StatusBadRequest)
		return
	}
	if err = png.Encode(&buf, img); err != nil {
		api.Error(w, r, fmt.Errorf("encode totp image failed: %w", err), http.StatusBadRequest)
		return
	}
	login := domain.Login{
		Key:      "",
		Provider: "totp",
		Subject:  otp.Secret(),
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
	data := map[string]interface{}{
		"loginId": login.Key,
		"secret":  otp.Secret(),
		"image":   "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
	}
	api.SuccessJson(w, r, data)
}

// EnableTOTP enables the totp login for the user
func (h *Handler) EnableTOTP(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	_ = urlParams.ByName("key") // user

	// TODO check if authenticated user is the same as the user in the url

	var totpEnableRequest domain.TotpEnableRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.Error(w, r, fmt.Errorf("read request body failed: %w", err), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &totpEnableRequest); err != nil {
		api.Error(w, r, fmt.Errorf("decode request body failed: %w", err), http.StatusBadRequest)
		return
	}
	login, err := h.db.Logins.Read(totpEnableRequest.LoginId, r.Context())
	if err != nil {
		api.Error(w, r, fmt.Errorf("read login failed: %w", err), http.StatusBadRequest)
		return
	}
	if !totp.Validate(totpEnableRequest.Code, login.Subject) {
		api.Error(w, r, fmt.Errorf("invalid totp code"), http.StatusBadRequest)
		return
	}
	if login.Enabled {
		api.Error(w, r, fmt.Errorf("totp already enabled"), http.StatusBadRequest)
		return
	}
	login.Enabled = true
	if err := h.db.Logins.Update(login, r.Context()); err != nil {
		api.Error(w, r, fmt.Errorf("update login failed: %w", err), http.StatusBadRequest)
		return
	}
	api.SuccessJson(w, r, nil)
}
