package user

import (
	"fmt"
	"pkv/api/src/repository/graph"
	"pkv/api/src/repository/security"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	db *graph.Db
}

func NewService(db *graph.Db) *Service {
	return &Service{db: db}
}

func UserToken(provider string, user string, expiry int64) string {
	return provider + ":" + user + ":" + strconv.FormatInt(expiry, 10)
}

func HashUserToken(token string) string {
	return security.HashToken(":user_authenticated::" + token)
}

func HashedUserToken(provider string, user string, expiry int64) string {
	token := UserToken(provider, user, expiry)
	return token + ":" + HashUserToken(token)
}

func ValidateUserToken(token string) (string, string, error) {
	// example: "x:username:expiry_unix:hash"
	tokens := strings.SplitN(token, ":", 4)
	if len(tokens) != 4 {
		return "", "", fmt.Errorf("token not correctly formatted")
	}
	unix := time.Now().Unix()
	expiry, err := strconv.ParseInt(tokens[2], 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("expiry not correctly formatted")
	}
	if unix > expiry {
		return "", "", fmt.Errorf("token expired")
	}
	hash := security.HashToken(":user_authenticated::" + tokens[0] + ":" + tokens[1] + ":" + tokens[2])
	if hash != tokens[3] {
		return "", "", fmt.Errorf("token invalid")
	}
	return tokens[1], tokens[0], nil
}
