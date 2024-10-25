package user

import (
	"context"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/t"
	"slices"
)

func (s *Service) Update(key string, name string, userType string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return t.Errorf("read user failed: %w", err)
	}
	if name == "" {
		name = user.Name
	}
	if len(name) > 100 {
		return t.Errorf("name cannot be longer than 100 characters")
	}
	if userType == "" {
		userType = user.Type
	}
	if !slices.Contains(dpv.ConfigInstance.Settings.UserTypes, userType) {
		return t.Errorf("invalid user type %v, choose one of the following: %+v", userType, dpv.ConfigInstance.Settings.UserTypes)
	}
	if userType == "administrator" && userType != user.Type {
		return t.Errorf("cannot update to administrator account")
	}
	user.Name = name
	user.Type = userType
	if err := s.db.Users.Update(user, ctx); err != nil {
		return t.Errorf("update user failed: %w", err)
	}
	return nil
}
