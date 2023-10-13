package user

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"pkv/api/src/domain"
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/user"
	"testing"
)

func TestExists(t *testing.T) {
	db, _, err := graph.Init("../../../config.yml", true)
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	var params httprouter.Params
	tests := []struct {
		name string
		code func(*graph.Db, http.HandlerFunc)
	}{
		{
			"fail if username is not valid",
			func(db *graph.Db, exists http.HandlerFunc) {
				params = httprouter.Params{{"key", "x"}}
				rr := callExists(exists, "/", t)
				if rr.Code != http.StatusBadRequest {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, http.StatusBadRequest)
				}
			},
		},
		{
			"exists returns false",
			func(db *graph.Db, exists http.HandlerFunc) {
				params = httprouter.Params{{"key", "doesnotexist"}}
				rr := callExists(exists, "/", t)
				if rr.Code != http.StatusOK {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, http.StatusOK)
				}
				var data bool
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				if err != nil {
					t.Fatalf("json unmarshalling failed: %s", err)
				}
				if data != false {
					t.Errorf("handler returned unexpected data: got %v want %v", data, false)
				}
			},
		},
		{
			"exists returns true",
			func(db *graph.Db, exists http.HandlerFunc) {
				user := domain.User{}
				err := db.Users.Create(&user, nil)
				if err != nil {
					t.Fatalf("user creation failed: %s", err)
				}
				params = httprouter.Params{{"key", user.Key}}
				rr := callExists(exists, "/", t)
				if rr.Code != http.StatusOK {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, http.StatusOK)
				}
				var data bool
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				if err != nil {
					t.Fatalf("json unmarshalling failed: %s", err)
				}
				if data != true {
					t.Errorf("handler returned unexpected data: got %v want %v", data, true)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := user.NewService(db)
			h := NewHandler(db, s)
			exists := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				h.Exists(writer, request, params)
			})
			tt.code(db, exists)
		})
	}
}

func callExists(exists http.HandlerFunc, url string, t *testing.T) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("request creation failed: %s", err)
	}
	rr := httptest.NewRecorder()
	exists.ServeHTTP(rr, req)
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
