package user

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"pkv/api/src/domain"
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/user"
	"testing"
	"time"
)

func TestPassword(t *testing.T) {
	db, _, err := graph.Init("../../../config.yml", true)
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	tests := []struct {
		name string
		code func(*graph.Db, http.HandlerFunc, http.HandlerFunc, *httprouter.Params)
	}{
		{
			"fail if user does not exist",
			func(db *graph.Db, linkPassword http.HandlerFunc, verifyPassword http.HandlerFunc, params *httprouter.Params) {
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatalf("request creation failed: %s", err)
				}
				loginAs("doesnotexist", req)
				rr := httptest.NewRecorder()
				*params = httprouter.Params{{"key", "doesnotexist"}}
				linkPassword.ServeHTTP(rr, req)
				expectedContentType := "application/json"
				if rr.Header().Get("Content-Type") != expectedContentType {
					t.Errorf("handler returned unexpected content-type: got %v want %v",
						rr.Header().Get("Content-Type"), expectedContentType)
				}
			},
		},
		{
			"happy path",
			func(db *graph.Db, linkPassword http.HandlerFunc, verifyPassword http.HandlerFunc, params *httprouter.Params) {
				user := domain.User{}
				err := db.Users.Create(&user, context.Background())
				if err != nil {
					t.Fatalf("user creation failed: %s", err)
				}
				*params = httprouter.Params{{"key", user.Key}}
				// attempt to set up no password
				rr := callLinkPassword(linkPassword, user.Key, "/", t)
				if rr.Code != 400 {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, 400)
				}
				// attempt to set up empty password
				rr = callLinkPassword(linkPassword, user.Key, "/?password=", t)
				if rr.Code != 400 {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, 400)
				}
				// attempt to set up simple password
				rr = callLinkPassword(linkPassword, user.Key, "/?password=123456", t)
				if rr.Code != 400 {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, 400)
				}
				// attempt to set up a normal password
				rr = callLinkPassword(linkPassword, user.Key, "/?password=Tr0ub4dor%263", t)
				if rr.Code != http.StatusOK {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, http.StatusOK)
				}
				// attempt to set up another password
				rr = callLinkPassword(linkPassword, user.Key, "/?password=Tr0ub4dor%264", t)
				if rr.Code != 400 {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, 400)
				}
				// attempt to verify the password
				rr = callVerifyPassword(verifyPassword, user.Key, "/?password=Tr0ub4dor%263", t)
				var data bool
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				if err != nil {
					t.Fatalf("json decoding failed: %s", err)
				}
				if !data {
					t.Errorf("handler returned unexpected data: got %v want %v", data, true)
				}
				// attempt to verify another password
				rr = callVerifyPassword(verifyPassword, user.Key, "/?password=Tr0ub4dor%264", t)
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				if err != nil {
					t.Fatalf("json decoding failed: %s", err)
				}
				if data {
					t.Errorf("handler returned unexpected data: got %v want %v", data, false)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params httprouter.Params
			t.Parallel()
			s := user.NewService(db)
			h := NewHandler(db, s)
			linkPassword := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				h.LinkPassword(writer, request, params)
			})
			verifyPassword := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				h.VerifyPassword(writer, request, params)
			})
			tt.code(db, linkPassword, verifyPassword, &params)
		})
	}
}

func callLinkPassword(linkPassword http.HandlerFunc, user, url string, t *testing.T) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", url, nil)
	loginAs(user, req)
	if err != nil {
		t.Fatalf("request creation failed: %s", err)
	}
	rr := httptest.NewRecorder()
	linkPassword.ServeHTTP(rr, req)
	expectedContentType := "application/json"
	// log.Printf("Status-Code: %d\n", rr.Code)
	// log.Printf("Content-Type: %s\n", rr.Header().Get("Content-Type"))
	// log.Printf("Body: %s\n", rr.Body.String())
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %v want %v",
			rr.Header().Get("Content-Type"), expectedContentType)
	}
	return rr
}

func callVerifyPassword(verifyPassword http.HandlerFunc, user, url string, t *testing.T) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("request creation failed: %s", err)
	}
	loginAs(user, req)
	rr := httptest.NewRecorder()
	verifyPassword.ServeHTTP(rr, req)
	expectedContentType := "application/json"
	//log.Printf("Status-Code: %d\n", rr.Code)
	//log.Printf("Content-Type: %s\n", rr.Header().Get("Content-Type"))
	//log.Printf("Body: %s\n", rr.Body.String())
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %v want %v",
			rr.Header().Get("Content-Type"), expectedContentType)
	}
	return rr
}

func loginAs(key string, r *http.Request) {
	r.Header.Set("Authorization", "user "+user.HashedUserToken("@", key, time.Now().Add(time.Minute*5).Unix()))
}
