package user

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pquerna/otp/totp"
	"image/png"
	"pkv/api/src/domain"
	"time"
)

// RequestTOTP generates a totp key and returns it as a base64 encoded png image
func (s *Service) RequestTOTP(key string, ctx context.Context) (map[string]interface{}, error) {
	// Check if the user exists
	_, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return nil, fmt.Errorf("read user failed: %w", err)
	}

	// Check if the user already has a TOTP login
	logins, err := s.db.GetLoginsForUser(key, ctx)
	if err != nil {
		return nil, fmt.Errorf("read logins failed: %w", err)
	}

	for _, login := range logins {
		if login.Provider == "totp" {
			return nil, fmt.Errorf("totp already requested")
		}
	}

	// Generate a TOTP key
	otp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "parkour-deutschland.de",
		AccountName: key, // TODO ensure that the authenticated user is the same as the user in the URL
	})
	if err != nil {
		return nil, fmt.Errorf("generate totp key failed: %w", err)
	}

	// Generate a base64-encoded PNG image
	var buf bytes.Buffer
	img, err := otp.Image(200, 200)
	if err != nil {
		return nil, fmt.Errorf("generate totp image failed: %w", err)
	}
	if err = png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode totp image failed: %w", err)
	}

	// Create a new TOTP login object
	login := domain.Login{
		Entity: domain.Entity{
			Key:     "",
			Created: time.Now(),
		},
		Provider: "totp",
		Subject:  otp.Secret(),
		Enabled:  false,
	}
	if err = s.db.Logins.Create(&login, ctx); err != nil {
		return nil, fmt.Errorf("create login failed: %w", err)
	}

	// Link the TOTP login to the user
	if err = s.db.LoginAuthenticatesUser(login, domain.User{Entity: domain.Entity{Key: key}}, ctx); err != nil {
		return nil, fmt.Errorf("link login to user failed: %w", err)
	}

	data := map[string]interface{}{
		"loginId": login.Key,
		"secret":  otp.Secret(),
		"image":   "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
	}

	return data, nil
}

// EnableTOTP enables the TOTP login for the user
func (s *Service) EnableTOTP(key string, totpEnableRequest domain.TotpEnableRequest, ctx context.Context) error {
	// TODO Check if the authenticated user is the same as the user in the URL
	// (Implement this part based on your authentication logic)

	// Read the login associated with the provided login ID
	login, err := s.db.Logins.Read(totpEnableRequest.LoginId, ctx)
	if err != nil {
		return fmt.Errorf("read login failed: %w", err)
	}

	// Validate the provided TOTP code
	if !totp.Validate(totpEnableRequest.Code, login.Subject) {
		return fmt.Errorf("invalid totp code")
	}

	// Check if TOTP is already enabled
	if login.Enabled {
		return fmt.Errorf("totp already enabled")
	}

	// Enable TOTP for the user
	login.Enabled = true
	if err := s.db.Logins.Update(login, ctx); err != nil {
		return fmt.Errorf("update login failed: %w", err)
	}

	return nil
}
