package authentication

import (
	"fmt"
	"pkv/api/src/internal/graph"
	"pkv/api/src/internal/security"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	db *graph.Db
}

func NewHandler(db *graph.Db) *Handler {
	return &Handler{db: db}
}

func hashUserToken(token string) string {
	return security.HashToken(":user_authenticated::" + token)
}

func ValidateUserToken(token string) error {
	// example: "x:username:expiry_unix:hash"
	tokens := strings.SplitN(token, ":", 4)
	if len(tokens) != 4 {
		return fmt.Errorf("token not correctly formatted")
	}
	unix := time.Now().Unix()
	expiry, err := strconv.ParseInt(tokens[2], 10, 64)
	if err != nil {
		return fmt.Errorf("expiry not correctly formatted")
	}
	if unix > expiry {
		return fmt.Errorf("token expired")
	}
	hash := security.HashToken(":user_authenticated::" + tokens[0] + ":" + tokens[1] + ":" + tokens[2])
	if hash != tokens[3] {
		return fmt.Errorf("token invalid")
	}
	return nil
}
