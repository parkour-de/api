package api

import (
	"encoding/json"
	"log"
	"net/http"
	"pkv/api/src/domain"
	"pkv/api/src/repository/graph"
	"pkv/api/src/repository/t"
	"pkv/api/src/service/user"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func CheckAuth(r *http.Request) (string, string, error) {
	auth := r.Header.Get("Authorization")
	format, auth, found := strings.Cut(auth, " ")
	if !found {
		return "", "", t.Errorf("authorization header missing")
	}
	if format != "user" {
		return "", "", t.Errorf("authorization header needs to start with 'user'")
	}
	key, method, err := user.ValidateUserToken(auth)
	if err != nil {
		return "", "", t.Errorf("invalid token: %w", err)
	}
	return key, method, nil
}

func Authenticated(r *http.Request) (string, error) {
	key, method, err := CheckAuth(r)
	if err != nil {
		return "", err
	}
	if method == "a" {
		return "", t.Errorf("your temporary login is expiring soon, please add a login method to your account first")
	}
	return key, nil
}

func IsAdmin(user domain.User) bool {
	return user.Type == "administrator"
}

func RequireSameUser(key string, r *http.Request) error {
	user, _, err := CheckAuth(r)
	if err != nil {
		return err
	}
	if key != user {
		return t.Errorf("you are logged in as %s, but you are trying to access %s", user, key)
	}
	return nil
}

func RequireUserAdmin(key string, r *http.Request, db *graph.Db) (string, string, error) {
	user, err := Authenticated(r)
	if key == "" {
		key = user
	}
	if err != nil {
		return "", "", err
	}
	if key != user {
		users, err := db.GetAdministeredUsers(user, r.Context())
		if err != nil {
			return "", "", t.Errorf("cannot get list of administered users: %w", err)
		}
		found := false
		for _, u := range users {
			if u.Key == key {
				found = true
				break
			}
		}
		if !found {
			return "", "", t.Errorf("user %s is not administered by %s", key, user)
		}
	}
	return key, user, nil
}

func RequireGlobalAdmin(r *http.Request, db *graph.Db) (*domain.User, error) {
	key, err := Authenticated(r)
	if err != nil {
		return nil, t.Errorf("authentication failed: %w", err)
	}
	user, err := db.Users.Read(key, r.Context())
	if err != nil {
		return nil, t.Errorf("reading current user failed: %w", err)
	}
	if !IsAdmin(*user) {
		return nil, t.Errorf("you are not an administrator")
	}
	return user, nil
}

func SuccessJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		Error(w, r, t.Errorf("serialising response failed: %w", err), 400)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		Success(w, r, jsonMsg)
	}
}

func Success(w http.ResponseWriter, r *http.Request, jsonMsg []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := w.Write(jsonMsg); err != nil {
		log.Printf("Error writing response: %v", err)
	}

	log.Printf(
		"%s %s %s 200",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
	)
}

func Error(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	if err == nil {
		err = t.Errorf("nil err")
	}
	logErr := err
	errorMsgJSON, err := json.Marshal(ErrorResponse{
		err.Error(),
	})
	if err != nil {
		log.Println(err)
	} else {
		if _, err = w.Write(errorMsgJSON); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}

	log.Printf(
		"%s %s %s %d %s",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		code,
		logErr.Error(),
	)
}

func MakeSet(queryParam string) map[string]struct{} {
	set := make(map[string]struct{})
	if queryParam != "" {
		tokens := strings.Split(queryParam, ",")
		for _, token := range tokens {
			set[token] = struct{}{}
		}
	}
	return set
}

func ParseInt(queryValue string) (int, error) {
	if queryValue == "" {
		return 0, nil
	}
	return strconv.Atoi(queryValue)
}

func ParseFloat(queryValue string) (float64, error) {
	if queryValue == "" {
		return 0, nil
	}
	return strconv.ParseFloat(queryValue, 64)
}
