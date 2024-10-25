package user

import (
	"context"
	"math/rand"
	"os"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/security"
	"pkv/api/src/repository/t"
	"strconv"
	"strings"
	"time"
)

func (s *Service) LinkPassword(key, password string, ctx context.Context) (string, error) {
	// Check password length
	if len(password) < 8 {
		return "", t.Errorf("password too short")
	}

	// Check password strength
	if !security.IsStrongPassword(password) {
		return "", t.Errorf("password too weak")
	}

	// Check if the user exists
	_, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return "", t.Errorf("read user failed: %w", err)
	}

	// Check if the user already has a password login
	logins, err := s.db.GetLoginsForUser(key, ctx)
	if err != nil {
		return "", t.Errorf("read logins failed: %w", err)
	}

	for _, login := range logins {
		if login.Provider == "password" {
			return "", t.Errorf("password already set")
		}
	}

	nonce := security.MakeNonce()
	login := domain.Login{
		Entity: domain.Entity{
			Key:     "",
			Created: time.Now(),
		},
		Provider: "password",
		Subject:  nonce + ":" + security.HashToken(":password::"+nonce+":"+password),
		Enabled:  true,
	}
	if err := s.db.Logins.Create(&login, ctx); err != nil {
		return "", t.Errorf("create login failed: %w", err)
	}

	if err := s.db.LoginAuthenticatesUser(login, domain.User{Entity: domain.Entity{Key: key}}, ctx); err != nil {
		return "", t.Errorf("link login to user failed: %w", err)
	}

	unix := time.Now().Unix()
	expiry := unix + 3600
	token := HashedUserToken("p", key, expiry)
	return token, nil
}

func (s *Service) Password(key, password string, ctx context.Context) (string, error) {
	success, err := s.VerifyPassword(key, password, ctx)
	if err != nil {
		return "", t.Errorf("verify password failed: %w", err)
	}
	if !success {
		return "", t.Errorf("password incorrect")
	}
	unix := time.Now().Unix()
	expiry := unix + 3600

	token := HashedUserToken("p", key, expiry)
	return token, nil
}

func (s *Service) VerifyPassword(key, password string, ctx context.Context) (bool, error) {
	// Check if the user exists
	_, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return false, t.Errorf("read user failed: %w", err)
	}

	// Check if the user has a password login
	logins, err := s.db.GetLoginsForUser(key, ctx)
	if err != nil {
		return false, t.Errorf("read logins failed: %w", err)
	}

	var login *domain.Login
	for _, l := range logins {
		if l.Provider == "password" {
			login = &l
		}
	}

	if login == nil {
		return false, t.Errorf("password login not found")
	}

	login, err = s.db.Logins.Read(login.Key, ctx)
	if err != nil {
		return false, t.Errorf("read login failed: %w", err)
	}

	if !login.Enabled {
		return false, t.Errorf("password login not enabled")
	}
	success := verifyPassword(password, login.Subject)

	return success, nil
}

var words3 string
var words4 string

func Suggest() (string, error) {
	sep := "-. /"
	randomSeparator := string(sep[rand.Intn(len(sep))])
	var err error
	if words3 == "" {
		words3, err = loadWords(dpv.ConfigInstance.Path + dpv.ConfigInstance.Server.Words1)
		if err != nil {
			return "", t.Errorf("load words failed: %w", err)
		}
	}
	words1list := strings.Split(words3, " ")
	if words4 == "" {
		words4, err = loadWords(dpv.ConfigInstance.Path + dpv.ConfigInstance.Server.Words2)
		if err != nil {
			return "", t.Errorf("load words failed: %w", err)
		}
	}
	words2list := strings.Split(words4, " ")
	var randomWords []string
	for i := 0; i < 5; i++ {
		randomWords = append(randomWords, words1list[rand.Intn(len(words1list))])
	}
	// replace 2 to 4 words with a random word from the second list
	amount := rand.Intn(3) + 2
	for i := 0; i < amount; i++ {
		randomWords[rand.Intn(len(randomWords))] = words2list[rand.Intn(len(words2list))]
	}
	return strings.Join(randomWords, randomSeparator), nil
}

func User() (string, error) {
	var err error
	if words3 == "" {
		words3, err = loadWords(dpv.ConfigInstance.Path + dpv.ConfigInstance.Server.Words3)
		if err != nil {
			return "", t.Errorf("load words failed: %w", err)
		}
	}
	words3list := strings.Split(words3, " ")
	if words4 == "" {
		words4, err = loadWords(dpv.ConfigInstance.Path + dpv.ConfigInstance.Server.Words4)
		if err != nil {
			return "", t.Errorf("load words failed: %w", err)
		}
	}
	words4list := strings.Split(words4, " ")
	return words3list[rand.Intn(len(words3list))] + words4list[rand.Intn(len(words4list))] + strconv.Itoa(rand.Intn(64)), nil
}

func loadWords(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		wd, _ := os.Getwd()
		return "", t.Errorf("could not load config file, looking for %v in %v: %w", filename, wd, err)
	}
	// file contains a string of words separated by spaces
	return string(bytes), nil
}

func verifyPassword(password, hash string) bool {
	parts := strings.Split(hash, ":")
	if len(parts) != 2 {
		return false
	}
	return security.HashToken(":password::"+parts[0]+":"+password) == parts[1]
}
