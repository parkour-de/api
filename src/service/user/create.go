package user

import (
	"context"
	"fmt"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
	"slices"
	"time"
)

func (s *Service) Create(key string, name string, userType string, ctx context.Context) (string, error) {
	if err := ValidateCustomKey(key); err != nil {
		return "", fmt.Errorf("invalid username: %w", err)
	}
	if name == "" {
		name = key
	}
	if len(name) > 100 {
		return "", fmt.Errorf("name cannot be longer than 100 characters")
	}
	if userType == "" {
		userType = dpv.ConfigInstance.Settings.UserTypes[0]
	}
	if !slices.Contains(dpv.ConfigInstance.Settings.UserTypes, userType) {
		return "", fmt.Errorf("invalid user type %v, choose one of the following: %+v", userType, dpv.ConfigInstance.Settings.UserTypes)
	}
	if userType == "administrator" {
		return "", fmt.Errorf("cannot create administrator account")
	}
	user := domain.User{
		Entity: domain.Entity{Key: key},
		Name:   name,
		Type:   userType,
		Information: map[string]string{
			"created": time.Now().Format(time.RFC3339),
			"login":   time.Now().Format(time.RFC3339),
		},
	}
	if err := s.db.Users.Create(&user, ctx); err != nil {
		return "", fmt.Errorf("create user failed: %w", err)
	}
	return user.Key, nil
}

func (s *Service) Claim(key string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return fmt.Errorf("read user failed: %w", err)
	}
	logins, err := s.db.GetLoginsForUser(key, ctx)
	if err != nil {
		return fmt.Errorf("read logins failed: %w", err)
	}
	if len(logins) > 0 {
		return fmt.Errorf("this username cannot be claimed")
	}
	administrators, err := s.db.GetAdministrators(key, ctx)
	if err != nil {
		return fmt.Errorf("read administrators failed: %w", err)
	}
	if len(administrators) > 0 {
		return fmt.Errorf("this username cannot be claimed")
	}
	creationDateString, ok := user.Information["created"]
	if !ok {
		return fmt.Errorf("user has no creation date")
	}
	creationDate, err := time.Parse(time.RFC3339, creationDateString)
	if err != nil {
		return fmt.Errorf("user has an invalid creation date")
	}
	if time.Now().Sub(creationDate) < 30*time.Minute {
		return fmt.Errorf("please wait %v more minutes before this username can be claimed", 30-time.Now().Sub(creationDate).Minutes())
	}
	err = s.db.Users.Delete(user, ctx)
	if err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}
	return nil
}
