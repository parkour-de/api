package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pkv/api/src/domain"
	"pkv/api/src/repository/graph"
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
		return "", "", fmt.Errorf("authorization header missing")
	}
	if format != "user" {
		return "", "", fmt.Errorf("authorization header needs to start with 'user'")
	}
	key, method, err := user.ValidateUserToken(auth)
	if err != nil {
		return "", "", fmt.Errorf("invalid token: %w", err)
	}
	return key, method, nil
}

func Authenticated(r *http.Request) (string, error) {
	key, method, err := CheckAuth(r)
	if err != nil {
		return "", err
	}
	if method == "a" {
		return "", fmt.Errorf("your temporary login is expiring soon, please add a login method to your account first")
	}
	return key, nil
}

func IsAdmin(user domain.User) bool {
	admin, ok := user.Information["admin"]
	return ok && admin == "true"
}

func RequireAdmin(r *http.Request, db *graph.Db) (*domain.User, error) {
	key, err := Authenticated(r)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	user, err := db.Users.Read(key, r.Context())
	if err != nil {
		return nil, fmt.Errorf("reading current user failed: %w", err)
	}
	if !IsAdmin(*user) {
		return nil, fmt.Errorf("you are not an administrator")
	}
	return user, nil
}

func SuccessJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		Error(w, r, fmt.Errorf("serialising response failed: %w", err), 400)
		return
	} else {
		Success(w, r, jsonMsg)
	}
}

func Success(w http.ResponseWriter, r *http.Request, jsonMsg []byte) {
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		err = fmt.Errorf("nil err")
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
