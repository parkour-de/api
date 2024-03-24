package user

import (
	"context"
	"fmt"
	"log"
	"pkv/api/src/domain"
	"pkv/api/src/repository/security"
	"time"
)

func (s *Service) RequestEmail(key string, email string, ctx context.Context) error {
	_, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return fmt.Errorf("read user failed: %w", err)
	}
	logins, err := s.db.GetLoginsForUser(key, ctx)
	if err != nil {
		return fmt.Errorf("read logins failed: %w", err)
	}
	for _, login := range logins {
		if login.Provider == "email" {
			return fmt.Errorf("email already requested")
		}
	}
	login := domain.Login{
		Entity: domain.Entity{
			Key:     "",
			Created: time.Now(),
		},
		Provider: "email",
		Subject:  email,
		Enabled:  false,
	}
	if err = s.db.Logins.Create(&login, ctx); err != nil {
		return fmt.Errorf("create login failed: %w", err)
	}
	if err = s.db.LoginAuthenticatesUser(login, domain.User{Entity: domain.Entity{Key: key}}, ctx); err != nil {
		return fmt.Errorf("link login to user failed: %w", err)
	}
	activationCode := emailActivationCode(login)
	activationLink := fmt.Sprintf("https://parkour-deutschland.de/user/%s/email/%s?code=%s", key, login.Key, activationCode)
	// TODO send email
	log.Println(activationLink)
	return nil
	// Clients: Send confirmation that an activation link has been sent.
	// Also users should make sure to check if the User id is correct. Clients should show the User id once more.
}

func emailActivationCode(login domain.Login) string {
	token := ":email_activation::" + login.Key + ":" + login.Subject
	activationCode := security.HashToken(token)
	return activationCode
}

func (s *Service) EnableEmail(loginId string, code string, ctx context.Context) error {
	login, err := s.db.Logins.Read(loginId, ctx)
	if err != nil {
		return fmt.Errorf("read login failed: %w", err)
	}
	if login.Provider != "email" {
		return fmt.Errorf("invalid provider")
	}
	if code != emailActivationCode(*login) {
		return fmt.Errorf("invalid activation code")
	}
	if login.Enabled {
		return fmt.Errorf("email already enabled")
	}
	login.Enabled = true
	if err := s.db.Logins.Update(login, ctx); err != nil {
		return fmt.Errorf("update login failed: %w", err)
	}
	return nil
}
