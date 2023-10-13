package graph

import (
	"pkv/api/src/domain"
	"testing"
)

func TestGetLoginsForUser(t *testing.T) {
	db, _, err := Init("../../../config.yml", true)
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	user := domain.User{}
	err = db.Users.Create(&user, nil)
	if err != nil {
		t.Fatalf("user creation failed: %s", err)
	}
	login1 := domain.Login{
		Provider: "NSA",
		Subject:  "do not reveal",
		Enabled:  true,
	}
	login2 := domain.Login{
		Provider: "CIA",
		Subject:  "do not reveal",
		Enabled:  false,
	}
	err = db.Logins.Create(&login1, nil)
	if err != nil {
		t.Fatalf("login creation failed: %s", err)
	}
	err = db.Logins.Create(&login2, nil)
	if err != nil {
		t.Fatalf("login creation failed: %s", err)
	}
	err = db.LoginAuthenticatesUser(login1, user, nil)
	if err != nil {
		t.Fatalf("linking login to user failed: %s", err)
	}
	err = db.LoginAuthenticatesUser(login2, user, nil)
	if err != nil {
		t.Fatalf("linking login to user failed: %s", err)
	}
	logins, err := db.GetLoginsForUser(user.Key, nil)
	if err != nil {
		t.Fatalf("get logins for user failed: %s", err)
	}
	if len(logins) != 2 {
		t.Fatalf("expected 2 logins, got %d", len(logins))
	}
	var nsaLogin, ciaLogin domain.Login
	for _, login := range logins {
		if login.Provider == "NSA" {
			nsaLogin = login
		} else if login.Provider == "CIA" {
			ciaLogin = login
		}
	}
	if nsaLogin.Key != login1.Key {
		t.Fatalf("expected login1, got %v", nsaLogin)
	}
	if ciaLogin.Key != login2.Key {
		t.Fatalf("expected login2, got %v", ciaLogin)
	}
	if nsaLogin.Enabled != login1.Enabled {
		t.Fatalf("expected login1 to be enabled, got %v", nsaLogin.Enabled)
	}
	if ciaLogin.Enabled != login2.Enabled {
		t.Fatalf("expected login2 to be disabled, got %v", ciaLogin.Enabled)
	}
	if nsaLogin.Subject != "" {
		t.Fatalf("expected login1 subject to be empty, got %v", nsaLogin.Subject)
	}
	if ciaLogin.Subject != "" {
		t.Fatalf("expected login2 subject to be empty, got %v", ciaLogin.Subject)
	}
}
