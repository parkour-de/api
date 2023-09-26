package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"net/url"
	"pkv/api/src/api"
	"pkv/api/src/internal/dpv"
	"pkv/api/src/internal/security"
	"strconv"
	"strings"
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

// Facebook handles the GET /api/facebook endpoint.
//
//	@Summary		Takes the authorization header from the browser and generates an access token for a user
//	@Description	Request an OAuth token from https://www.facebook.com/v17.0/dialog/oauth, then call this endpoint
//	@Description	using Authorization: facebook mySuperSecretToken as a header. This endpoint will then make a debug
//	@Description	call to Facebook Graph API to extract a unique ID that can be attached to a user.
//	@Tags			authentication
//	@Success		200	{object}	string
//	@Router			/api/facebook [get]
func (h *Handler) Facebook(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	if api.MakeCors(w, r) {
		return
	}
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "facebook ") {
		tokens := strings.SplitN(auth, " ", 2)
		if len(tokens) != 2 {
			api.Error(w, fmt.Errorf("authorization header not correctly formatted"), http.StatusBadRequest)
			return
		}
		auth = tokens[1]
		if len(auth) < 1 {
			api.Error(w, fmt.Errorf("authorization header contains empty token"), http.StatusBadRequest)
			return
		}
	} else {
		api.Error(w, fmt.Errorf("authorization header missing or not correctly prefixed"), http.StatusBadRequest)
		return
	}

	validationURL := dpv.ConfigInstance.Auth.FacebookGraphUrl
	params := url.Values{}
	params.Set("input_token", auth)
	params.Set("access_token", auth)

	resp, err := http.Get(validationURL + "?" + params.Encode())
	if err != nil {
		fmt.Printf("token validation failed - %s\n", err.Error())
		api.Error(w, fmt.Errorf("token validation failed - check server logs"), http.StatusBadRequest)
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
		api.Error(w, fmt.Errorf("token validation failed - check server logs"), http.StatusUnauthorized)
		return
	}

	// Parse the validation response
	var validationResponse FacebookTokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !validationResponse.Data.IsValid {
		api.Error(w, fmt.Errorf("facebook says, data is not valid"), http.StatusUnauthorized)
		return
	}

	if validationResponse.Data.AppId != dpv.ConfigInstance.Auth.FacebookAppId {
		api.Error(w, fmt.Errorf("facebook says, this token belongs to a different app"), http.StatusUnauthorized)
		return
	}

	now := time.Now()
	iat := int64(validationResponse.Data.IssuedAt)
	exp := int64(validationResponse.Data.ExpiresAt)
	unix := now.Unix()

	if iat > unix {
		api.Error(w, fmt.Errorf("facebook says, this token is from the future"), http.StatusUnauthorized)
		return
	}

	if exp < unix {
		api.Error(w, fmt.Errorf("facebook says, this token is from the past"), http.StatusUnauthorized)
		return
	}

	expiry := unix + 3600
	if exp < expiry {
		expiry = exp
	}

	user := "4711"

	token := user + "." + strconv.FormatInt(expiry, 10)

	hash := security.HashToken(token)

	token = token + "." + hash

	jsonMsg, err := json.Marshal(token)
	if err != nil {
		api.Error(w, fmt.Errorf("converting token failed: %w", err), 400)
		return
	}

	api.Success(w, jsonMsg)
}
