package user

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"net/url"
	"pkv/api/src/api"
	"pkv/api/src/domain"
	"pkv/api/src/internal/dpv"
	"time"
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

// LinkFacebook attaches a facebook subject (sub) to a user
func (h *Handler) LinkFacebook(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	key := urlParams.ByName("key")
	user, err := h.db.Users.Read(key, r.Context())
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
		if login.Provider == "facebook" {
			api.Error(w, r, fmt.Errorf("facebook already connected"), http.StatusBadRequest)
			return
		}
	}
	auth := r.URL.Query().Get("auth")

	validationURL := dpv.ConfigInstance.Auth.FacebookGraphUrl
	params := url.Values{}
	params.Set("input_token", auth)
	params.Set("access_token", auth)

	resp, err := http.Get(validationURL + "?" + params.Encode())
	if err != nil {
		fmt.Printf("token validation failed - %s\n", err.Error())
		api.Error(w, r, fmt.Errorf("token validation failed - check server logs"), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("token validation failed - %d\n", resp.StatusCode)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("could not read associated error message - %s\n", err.Error())
		} else {
			fmt.Println(string(bodyBytes))
		}
		api.Error(w, r, fmt.Errorf("token validation failed - check server logs"), http.StatusUnauthorized)
		return
	}

	// Parse the validation response
	var validationResponse FacebookTokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !validationResponse.Data.IsValid {
		api.Error(w, r, fmt.Errorf("facebook says, data is not valid"), http.StatusUnauthorized)
		return
	}

	if validationResponse.Data.AppId != dpv.ConfigInstance.Auth.FacebookAppId {
		api.Error(w, r, fmt.Errorf("facebook says, this token belongs to a different app"), http.StatusUnauthorized)
		return
	}

	now := time.Now()
	iat := int64(validationResponse.Data.IssuedAt)
	exp := int64(validationResponse.Data.ExpiresAt)
	unix := now.Unix()

	if iat > unix {
		api.Error(w, r, fmt.Errorf("facebook says, this token is from the future"), http.StatusUnauthorized)
		return
	}

	if exp < unix {
		api.Error(w, r, fmt.Errorf("facebook says, this token is from the past"), http.StatusUnauthorized)
		return
	}

	expiry := unix + 3600
	if exp < expiry {
		expiry = exp
	}

	login := domain.Login{
		Key:      "",
		Provider: "facebook",
		Subject:  validationResponse.Data.UserId,
		Enabled:  true,
		Created:  time.Now(),
	}

	if err = h.db.Logins.Create(&login, r.Context()); err != nil {
		api.Error(w, r, fmt.Errorf("create login failed: %w", err), http.StatusBadRequest)
		return
	}

	if err = h.db.LoginAuthenticatesUser(login, *user, r.Context()); err != nil {
		api.Error(w, r, fmt.Errorf("link login to user failed: %w", err), http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, nil)
}
