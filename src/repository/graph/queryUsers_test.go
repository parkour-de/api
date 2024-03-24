package graph

import (
	"context"
	"pkv/api/src/domain"
	"slices"
	"testing"
)

func TestGetAssumableUsers(t *testing.T) {
	db, _, err := Init("../../../config.yml", true)
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	users := []domain.User{
		{Entity: domain.Entity{Key: "admin"}},
		{Entity: domain.Entity{Key: "u1"}},
		{Entity: domain.Entity{Key: "u11"}},
		{Entity: domain.Entity{Key: "u12"}},
		{Entity: domain.Entity{Key: "decoy1"}},
		{Entity: domain.Entity{Key: "decoy2"}},
		{Entity: domain.Entity{Key: "decoy3"}},
	}
	for _, user := range users {
		err = db.Users.Create(&user, context.Background())
		if err != nil {
			t.Fatalf("user creation failed: %s", err)
		}
	}
	err = db.UserAdministersUser(users[0], users[1], context.Background())
	if err != nil {
		t.Fatalf("linking user to user failed: %s", err)
	}
	err = db.UserAdministersUser(users[1], users[2], context.Background())
	if err != nil {
		t.Fatalf("linking user to user failed: %s", err)
	}
	err = db.UserAdministersUser(users[1], users[3], context.Background())
	if err != nil {
		t.Fatalf("linking user to user failed: %s", err)
	}
	err = db.UserAdministersUser(users[2], users[3], context.Background())
	if err != nil {
		t.Fatalf("linking user to user failed: %s", err)
	}
	if _, err := db.Edges.CreateDocument(context.Background(), domain.Edge{
		From:  "users/" + users[0].Key,
		To:    "users/" + users[4].Key,
		Label: "decoys",
	}); err != nil {
		t.Fatalf("Link failed: %s", err)
	}
	if _, err := db.Edges.CreateDocument(context.Background(), domain.Edge{
		From:  "users/" + users[4].Key,
		To:    "users/" + users[5].Key,
		Label: "decoys",
	}); err != nil {
		t.Fatalf("Link failed: %s", err)
	}
	if _, err := db.Edges.CreateDocument(context.Background(), domain.Edge{
		From:  "users/" + users[4].Key,
		To:    "users/" + users[6].Key,
		Label: "administers",
	}); err != nil {
		t.Fatalf("Link failed: %s", err)
	}

	assumableUsers, err := db.GetAdministeredUsers(users[0].Key, context.Background())
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
		if !slices.Contains(assumableUserKeys, user.Key) {
			t.Fatalf("expected %s to be assumable, got %v", user.Key, assumableUserKeys)
		}
	}
}
