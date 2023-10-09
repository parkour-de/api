package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"pkv/api/src/domain"
	"pkv/api/src/internal/dpv"
	"strconv"
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

func (s *Service) LinkFacebook(key string, auth string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return fmt.Errorf("read user failed: %w", err)
	}
	logins, err := s.db.GetLoginsForUser(key, ctx)
	if err != nil {
		return fmt.Errorf("read logins failed: %w", err)
	}
	for _, login := range logins {
		if login.Provider == "facebook" {
			return fmt.Errorf("facebook already connected")
		}
	}

	validationResponse, err := s.checkFacebookAuth(auth)
	if err != nil {
		return err
	}

	login := domain.Login{
		Key:      "",
		Provider: "facebook",
		Subject:  validationResponse.Data.UserId,
		Enabled:  true,
		Created:  time.Now(),
	}

	if err = s.db.Logins.Create(&login, ctx); err != nil {
		return fmt.Errorf("create login failed: %w", err)
	}

	if err = s.db.LoginAuthenticatesUser(login, *user, ctx); err != nil {
		return fmt.Errorf("link login to user failed: %w", err)
	}
	return nil
}

func (s *Service) Facebook(auth string) (string, error) {
	validationResponse, err := s.checkFacebookAuth(auth)
	if err != nil {
		return "", err
	}
	exp := int64(validationResponse.Data.ExpiresAt)
	unix := time.Now().Unix()

	expiry := unix + 3600
	if exp < expiry {
		expiry = exp
	}

	user := "4711"

	token := facebookToken(user, expiry)

	hash := hashUserToken(token)

	token = token + ":" + hash
	return token, nil
}

func facebookToken(user string, expiry int64) string {
	return "f:" + user + ":" + strconv.FormatInt(expiry, 10)
}

func (s *Service) checkFacebookAuth(auth string) (FacebookTokenValidationResponse, error) {
	validationURL := dpv.ConfigInstance.Auth.FacebookGraphUrl
	params := url.Values{}
	params.Set("input_token", auth)
	params.Set("access_token", auth)

	resp, err := http.Get(validationURL + "?" + params.Encode())
	if err != nil {
		log.Printf("token validation failed - %s\n", err.Error())
		return FacebookTokenValidationResponse{}, fmt.Errorf("token validation failed - check server logs")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("token validation failed - %d\n", resp.StatusCode)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("could not read associated error message - %s\n", err.Error())
		} else {
			log.Println(string(bodyBytes))
		}
		return FacebookTokenValidationResponse{}, fmt.Errorf("token validation failed - check server logs")
	}

	// Parse the validation response
	var validationResponse FacebookTokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
		log.Printf("could not parse validation response - %s\n", err.Error())
		return FacebookTokenValidationResponse{}, fmt.Errorf("could not parse validation response - check server logs")
	}

	if !validationResponse.Data.IsValid {
		return FacebookTokenValidationResponse{}, fmt.Errorf("facebook says, data is not valid")
	}

	if validationResponse.Data.AppId != dpv.ConfigInstance.Auth.FacebookAppId {
		return FacebookTokenValidationResponse{}, fmt.Errorf("facebook says, this token belongs to a different app")
	}

	iat := int64(validationResponse.Data.IssuedAt)
	exp := int64(validationResponse.Data.ExpiresAt)
	unix := time.Now().Unix()

	if iat > unix {
		return FacebookTokenValidationResponse{}, fmt.Errorf("facebook says, this token is from the future")
	}

	if exp < unix {
		return FacebookTokenValidationResponse{}, fmt.Errorf("facebook says, this token is from the past")
	}
	return validationResponse, nil
}
