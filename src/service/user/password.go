package user

import (
	"context"
	"fmt"
	"log"
	"pkv/api/src/domain"
	"pkv/api/src/repository/security"
	"strings"
	"time"
)

func (s *Service) LinkPassword(key, password string, ctx context.Context) error {
	// Check password length
	if len(password) < 8 {
		return fmt.Errorf("password too short")
	}

	// Check password strength
	if !security.IsStrongPassword(password) {
		return fmt.Errorf("password too weak")
	}

	// Check if the user exists
	_, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return fmt.Errorf("read user failed: %w", err)
	}

	// Check if the user already has a password login
	logins, err := s.db.GetLoginsForUser(key, ctx)
	if err != nil {
		return fmt.Errorf("read logins failed: %w", err)
	}

	for _, login := range logins {
		if login.Provider == "password" {
			return fmt.Errorf("password already set")
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
	log.Println("pwd: " + password)
	log.Println("sub: " + login.Subject)
	if err := s.db.Logins.Create(&login, ctx); err != nil {
		return fmt.Errorf("create login failed: %w", err)
	}

	if err := s.db.LoginAuthenticatesUser(login, domain.User{Key: key}, ctx); err != nil {
		return fmt.Errorf("link login to user failed: %w", err)
	}

	return nil
}

func (s *Service) VerifyPassword(key, password string, ctx context.Context) (bool, error) {
	// Check if the user exists
	_, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return false, fmt.Errorf("read user failed: %w", err)
	}

	// Check if the user has a password login
	logins, err := s.db.GetLoginsForUser(key, ctx)
	if err != nil {
		return false, fmt.Errorf("read logins failed: %w", err)
	}

	var login *domain.Login
	for _, l := range logins {
		if l.Provider == "password" {
			login = &l
		}
	}

	if login == nil {
		return false, fmt.Errorf("password login not found")
	}

	login, err = s.db.Logins.Read(login.Key, ctx)
	if err != nil {
		return false, fmt.Errorf("read login failed: %w", err)
	}

	if !login.Enabled {
		return false, fmt.Errorf("password login not enabled")
	}

	log.Println("pwd: " + password)
	log.Println("sub: " + login.Subject)
	success := verifyPassword(password, login.Subject)

	return success, nil
}

func verifyPassword(password, hash string) bool {
	parts := strings.Split(hash, ":")
	if len(parts) != 2 {
		return false
	}
	return security.HashToken(":password::"+parts[0]+":"+password) == parts[1]
}
