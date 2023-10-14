package user

import (
	"context"
	"fmt"
	"regexp"
)

func (s *Service) Exists(username string, ctx context.Context) (bool, error) {
	if err := ValidateKey(username); err != nil {
		return false, fmt.Errorf("invalid username: %w", err)
	}

	// Check if the user exists
	exists, err := s.db.Users.Has(username, ctx)
	if err != nil {
		return false, fmt.Errorf("check user exists failed: %w", err)
	}

	return exists, nil
}

func ValidateCustomKey(username string) error {
	if matched, _ := regexp.MatchString(`^\d+$`, username); matched {
		return fmt.Errorf("key cannot only contain digits")
	}
	return ValidateKey(username)
}

func ValidateKey(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters long")
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-][a-zA-Z0-9_\-.]{2,29}$`, username); !matched {
		return fmt.Errorf("key must contain a-z, A-Z, 0-9, _, -, or . but may not start with a period")
	}
	return nil
}
