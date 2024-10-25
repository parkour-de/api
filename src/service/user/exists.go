package user

import (
	"context"
	"pkv/api/src/repository/t"
	"regexp"
)

func (s *Service) Exists(username string, ctx context.Context) (bool, error) {
	if err := ValidateKey(username); err != nil {
		return false, t.Errorf("invalid username: %w", err)
	}

	// Check if the user exists
	exists, err := s.db.Users.Has(username, ctx)
	if err != nil {
		return false, t.Errorf("check user exists failed: %w", err)
	}

	return exists, nil
}

func ValidateCustomKey(username string) error {
	if matched, _ := regexp.MatchString(`^\d+$`, username); matched {
		return t.Errorf("key cannot only contain digits")
	}
	return ValidateKey(username)
}

func ValidateKey(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return t.Errorf("username must be between 3 and 30 characters long")
	}
	if matched, _ := regexp.MatchString(`^[a-z0-9_-][a-z0-9_\-.]{2,29}$`, username); !matched {
		return t.Errorf("key must contain a-z, 0-9, _, -, or . but may not start with a period")
	}
	return nil
}
