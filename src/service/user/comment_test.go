package user

import (
	"pkv/api/src/domain"
	"pkv/api/src/repository/graph"
	"testing"
)

func TestComments(t *testing.T) {
	db, _, err := graph.Init("../../../config.yml", true)
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	user := domain.User{}
	err = db.Users.Create(&user, nil)
	if err != nil {
		t.Fatalf("initialisation failed: %s", err)
	}
	service := NewService(db)
	err = service.AddComment(user.Key, "author", "title", "text", nil)
	if err != nil {
		t.Fatalf("add comment failed: %s", err)
	}
	err = service.AddComment(user.Key, "author", "title2", "text2", nil)
	if err != nil {
		t.Fatalf("add comment failed: %s", err)
	}
	err = service.AddComment(user.Key, "author", "title2", "text2", nil)
	if err == nil {
		t.Fatalf("add comment should fail")
	}
	err = service.EditComment(user.Key, "author", "title2", "title3", "text3", nil)
	if err != nil {
		t.Fatalf("edit comment failed: %s", err)
	}
	err = service.EditComment(user.Key, "author", "title2", "title4", "text4", nil)
	if err == nil {
		t.Fatalf("edit comment should fail")
	}
	err = service.EditComment(user.Key, "author", "title3", "title3", "text3", nil)
	if err != nil {
		t.Fatalf("edit comment failed: %s", err)
	}
	err = service.EditComment(user.Key, "wrong_user", "title3", "title3", "text3", nil)
	if err == nil {
		t.Fatalf("edit comment should fail")
	}
	puser, err := db.Users.Read(user.Key, nil)
	if err != nil {
		t.Fatalf("read user failed: %s", err)
	}
	if len(puser.Comments) != 2 {
		t.Fatalf("wrong number of comments: %d", len(puser.Comments))
	}
	if puser.Comments[1].Title != "title3" {
		t.Fatalf("wrong comment title: %s", puser.Comments[0].Title)
	}
	if puser.Comments[1].Text != "text3" {
		t.Fatalf("wrong comment text: %s", puser.Comments[0].Text)
	}
	if puser.Comments[1].Author != "author" {
		t.Fatalf("wrong comment author: %s", puser.Comments[0].Author)
	}
	err = service.DeleteComment(user.Key, "wrong_user", "title3", nil)
	if err == nil {
		t.Fatalf("delete comment should fail")
	}
	err = service.DeleteComment(user.Key, "author", "title3", nil)
	if err != nil {
		t.Fatalf("delete comment failed: %s", err)
	}
	puser, err = db.Users.Read(user.Key, nil)
	if err != nil {
		t.Fatalf("read user failed: %s", err)
	}
	if len(puser.Comments) != 1 {
		t.Fatalf("wrong number of comments: %d", len(puser.Comments))
	}
	if puser.Comments[0].Title != "title" {
		t.Fatalf("wrong comment title: %s", puser.Comments[0].Title)
	}
	if puser.Comments[0].Text != "text" {
		t.Fatalf("wrong comment text: %s", puser.Comments[0].Text)
	}
}
