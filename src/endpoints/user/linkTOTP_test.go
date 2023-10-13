package user

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/pquerna/otp/totp"
	"log"
	"net/http"
	"net/http/httptest"
	"pkv/api/src/domain"
	"pkv/api/src/repository/graph"
	"pkv/api/src/service/user"
	"testing"
	"time"
)

func TestTOTP(t *testing.T) {
	db, _, err := graph.Init("../../../config.yml", true)
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	var params httprouter.Params

	tests := []struct {
		name  string
		setup func() *graph.Db
		code  func(*graph.Db, http.HandlerFunc, http.HandlerFunc)
	}{
		{
			"fail if user does not exist",
			func() *graph.Db {
				return db
			},
			func(db *graph.Db, requestTOTP http.HandlerFunc, enableTOTP http.HandlerFunc) {
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatalf("request creation failed: %s", err)
				}
				rr := httptest.NewRecorder()
				params = httprouter.Params{{"key", "doesnotexist"}}
				requestTOTP.ServeHTTP(rr, req)
				expectedContentType := "application/json"
				if rr.Header().Get("Content-Type") != expectedContentType {
					t.Errorf("handler returned unexpected content-type: got %v want %v",
						rr.Header().Get("Content-Type"), expectedContentType)
				}
				// TODO: Actually fail the test
				if rr.Code != http.StatusBadRequest {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, http.StatusBadRequest)
				}
			},
		},
		{
			"happy path",
			func() *graph.Db {
				return db
			},
			func(db *graph.Db, requestTOTP http.HandlerFunc, enableTOTP http.HandlerFunc) {
				user := domain.User{}
				err := db.Users.Create(&user, nil)
				if err != nil {
					t.Fatalf("user creation failed: %s", err)
				}
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatalf("request creation failed: %s", err)
				}
				rr := httptest.NewRecorder()
				params = httprouter.Params{{"key", user.Key}}
				requestTOTP.ServeHTTP(rr, req)
				expectedContentType := "application/json"
				if rr.Code != http.StatusOK {
					t.Errorf("handler returned unexpected status code: got %v want %v", rr.Code, http.StatusOK)
				}
				if rr.Header().Get("Content-Type") != expectedContentType {
					t.Errorf("handler returned unexpected content-type: got %v want %v",
						rr.Header().Get("Content-Type"), expectedContentType)
				}
				data := map[string]string{}
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				if err != nil {
					t.Fatalf("could not decode response body: %s", err)
				}
				loginId, ok := data["loginId"]
				if !ok {
					t.Fatalf("loginId not found in response body")
				}
				secret, ok := data["secret"]
				if !ok {
					t.Fatalf("secret not found in response body")
				}

				login, err := db.Logins.Read(loginId, nil)
				if err != nil {
					t.Fatalf("could not read login: %s", err)
				}
				if login.Subject != secret {
					t.Fatalf("secret does not match")
				}
				if login.Enabled {
					t.Fatalf("login should not be enabled")
				}
				if login.Provider != "totp" {
					t.Fatalf("provider should be totp")
				}

				wrongCode := "1234567"

				enableRequest := domain.TotpEnableRequest{
					LoginId: loginId,
					Code:    wrongCode,
				}
				body, err := json.Marshal(enableRequest)
				if err != nil {
					t.Fatalf("could not encode request body: %s", err)
				}
				req, err = http.NewRequest("POST", "/", bytes.NewReader(body))
				if err != nil {
					t.Fatalf("request creation failed: %s", err)
				}
				rr = httptest.NewRecorder()
				enableTOTP.ServeHTTP(rr, req)
				log.Printf("Status-Code: %d\n", rr.Code)
				if rr.Code != http.StatusBadRequest {
					t.Errorf("with the wrong code, handler returned unexpected status code: got %v want %v", rr.Code, http.StatusBadRequest)
				}

				code, err := totp.GenerateCode(secret, time.Now())
				if err != nil {
					t.Fatalf("could not generate code: %s", err)
				}
				enableRequest = domain.TotpEnableRequest{
					LoginId: loginId,
					Code:    code,
				}
				body, err = json.Marshal(enableRequest)
				if err != nil {
					t.Fatalf("could not encode request body: %s", err)
				}
				req, err = http.NewRequest("POST", "/", bytes.NewReader(body))
				if err != nil {
					t.Fatalf("request creation failed: %s", err)
				}
				rr = httptest.NewRecorder()
				enableTOTP.ServeHTTP(rr, req)
				log.Printf("Status-Code: %d\n", rr.Code)
				if rr.Code != http.StatusOK {
					t.Errorf("with the correct code, handler returned unexpected status code: got %v want %v", rr.Code, http.StatusOK)
				}

				login, err = db.Logins.Read(loginId, nil)
				if err != nil {
					t.Fatalf("could not read login: %s", err)
				}
				if !login.Enabled {
					t.Fatalf("login should be enabled")
				}

				req, err = http.NewRequest("POST", "/", bytes.NewReader(body))
				if err != nil {
					t.Fatalf("request creation failed: %s", err)
				}
				rr = httptest.NewRecorder()
				enableTOTP.ServeHTTP(rr, req)
				log.Printf("Status-Code: %d\n", rr.Code)
				if rr.Code != http.StatusBadRequest {
					t.Errorf("with the repeated activation, handler returned unexpected status code: got %v want %v", rr.Code, http.StatusBadRequest)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params = httprouter.Params{}
			db := tt.setup()
			s := user.NewService(db)
			h := NewHandler(db, s)
			requestTOTP := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				h.RequestTOTP(writer, request, params)
			})
			enableTOTP := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				h.EnableTOTP(writer, request, params)
			})
			tt.code(db, requestTOTP, enableTOTP)
		})
	}
}
