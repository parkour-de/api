package graph

import (
	"pkv/api/src/domain"
	"testing"
)

func TestGetAssumableUsers(t *testing.T) {
	db, _, err := Init("../../../config.yml", true)
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	users := []domain.User{
		{Key: "admin"},
		{Key: "u1"},
		{Key: "u11"},
		{Key: "u12"},
		{Key: "decoy1"},
		{Key: "decoy2"},
		{Key: "decoy3"},
	}
	for _, user := range users {
		err = db.Users.Create(&user, nil)
		if err != nil {
			t.Fatalf("user creation failed: %s", err)
		}
	}
	err = db.UserAdministersUser(users[0], users[1], nil)
	if err != nil {
		t.Fatalf("linking user to user failed: %s", err)
	}
	err = db.UserAdministersUser(users[1], users[2], nil)
	if err != nil {
		t.Fatalf("linking user to user failed: %s", err)
	}
	err = db.UserAdministersUser(users[1], users[3], nil)
	if err != nil {
		t.Fatalf("linking user to user failed: %s", err)
	}
	err = db.UserAdministersUser(users[2], users[3], nil)
	if err != nil {
		t.Fatalf("linking user to user failed: %s", err)
	}
	if _, err := db.Edges.CreateDocument(nil, domain.Edge{
		From:  "users/" + users[0].Key,
		To:    "users/" + users[4].Key,
		Label: "decoys",
	}); err != nil {
		t.Fatalf("Link failed: %s", err)
	}
	if _, err := db.Edges.CreateDocument(nil, domain.Edge{
		From:  "users/" + users[4].Key,
		To:    "users/" + users[5].Key,
		Label: "decoys",
	}); err != nil {
		t.Fatalf("Link failed: %s", err)
	}
	if _, err := db.Edges.CreateDocument(nil, domain.Edge{
		From:  "users/" + users[4].Key,
		To:    "users/" + users[6].Key,
		Label: "administers",
	}); err != nil {
		t.Fatalf("Link failed: %s", err)
	}

	assumableUsers, err := db.GetAdministeredUsers(users[0].Key, nil)
	if err != nil {
		t.Fatalf("get assumable users failed: %s", err)
	}
	var assumableUserKeys []string
	for _, user := range assumableUsers {
		assumableUserKeys = append(assumableUserKeys, user.Key)
	}
	if len(assumableUserKeys) != 4 {
		t.Fatalf("expected 4 users, got %v", assumableUserKeys)
	}
	// compare with users[0:3] ignoring order
	for _, user := range users[0:3] {
		if !contains(assumableUserKeys, user.Key) {
			t.Fatalf("expected %s to be assumable, got %v", user.Key, assumableUserKeys)
		}
	}
}
func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
