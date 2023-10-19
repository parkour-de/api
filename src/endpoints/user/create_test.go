package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/user"
	"testing"
)

func TestCreate(t *testing.T) {
	db, config, err := graph.Init("../../../config.yml", true)
	dpv.ConfigInstance = config
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	service := user.NewService(db)
	handler := NewHandler(db, service)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("request creation failed: %s", err)
	}
	rr := httptest.NewRecorder()
	handler.Create(rr, req, httprouter.Params{{"key", "x"}}) // too short
	if rr.Code != 400 {
		t.Errorf("should have rejected username that is too short: got %v want %v", rr.Code, 400)
	}
	rr = httptest.NewRecorder()
	handler.Create(rr, req, httprouter.Params{{"key", "4"}}) // only numbers
	if rr.Code != 400 {
		t.Errorf("should have rejected username that is only numbers: got %v want %v", rr.Code, 400)
	}
	rr = httptest.NewRecorder()
	handler.Create(rr, req, httprouter.Params{{"key", "MiXeD"}}) // mixed case
	if rr.Code != 400 {
		t.Errorf("should have rejected username that contains upper case characters: got %v want %v", rr.Code, 400)
	}
	rr = httptest.NewRecorder()
	handler.Create(rr, req, httprouter.Params{{"key", "hello"}}) // ok
	if rr.Code != 200 {
		t.Errorf("should have accepted username: got %v want %v", rr.Code, 200)
	}
	body := rr.Body.String()
	t.Logf("body: %s", body)
	key, method, err := user.ValidateUserToken(body)
	if err != nil {
		t.Errorf("should have returned a valid token: %s", err)
	}
	if method != "a" {
		t.Errorf("should have returned a token with method 'a': got %v want %v", method, "a")
	}
	if key != "hello" {
		t.Errorf("should have returned a token with key 'hello': got %v want %v", key, "hello")
	}
}

func TestClaim(t *testing.T) {
	db, config, err := graph.Init("../../../config.yml", true)
	dpv.ConfigInstance = config
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	service := user.NewService(db)
	handler := NewHandler(db, service)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("request creation failed: %s", err)
	}
	rr := httptest.NewRecorder()
	handler.Create(rr, req, httprouter.Params{{"key", "hello"}})
	if rr.Code != 200 {
		t.Errorf("should have accepted username: got %v want %v", rr.Code, 200)
	}
	rr = httptest.NewRecorder()
	handler.Claim(rr, req, httprouter.Params{{"key", "hello"}}) // taken
	if rr.Code != 400 {
		t.Errorf("should have rejected username that is taken: got %v want %v", rr.Code, 400)
	}
	user, err := db.Users.Read("hello", nil)
	if err != nil {
		t.Fatalf("user read failed: %s", err)
	}
	user.Information["created"] = "1970-01-05T15:04:05Z"
	if err := db.Users.Update(user, nil); err != nil {
		t.Fatalf("user update failed: %s", err)
	}
	rr = httptest.NewRecorder()
	handler.Claim(rr, req, httprouter.Params{{"key", "hello"}}) // release
	if rr.Code != 200 {
		t.Errorf("should have deleted username that is taken: got %v want %v", rr.Code, 200)
	}
	_, err = db.Users.Read("hello", nil)
	if err == nil {
		t.Errorf("should have deleted username that is taken")
	}
}
