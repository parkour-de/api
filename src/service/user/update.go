package user

import (
	"context"
	"fmt"
	"pkv/api/src/repository/dpv"
	"slices"
)

func (s *Service) Update(key string, name string, userType string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return fmt.Errorf("read user failed: %w", err)
	}
	if name == "" {
		name = user.Name
	}
	if len(name) > 100 {
		return fmt.Errorf("name cannot be longer than 100 characters")
	}
	if userType == "" {
		userType = user.Type
	}
	if !slices.Contains(dpv.ConfigInstance.Settings.UserTypes, userType) {
		return fmt.Errorf("invalid user type %v, choose one of the following: %+v", userType, dpv.ConfigInstance.Settings.UserTypes)
	}
	user.Name = name
	user.Type = userType
	if err := s.db.Users.Update(user, ctx); err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}
