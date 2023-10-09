package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/sha3"
	"pkv/api/src/internal/dpv"
)

func MakeNonce() string {
	buff := make([]byte, 12)
	_, err := rand.Read(buff)
	if err != nil {
		println(fmt.Errorf("random number generation failed: %w", err).Error())
	}
	return base64.RawURLEncoding.EncodeToString(buff)
}

func HashToken(token string) string {
	secretKey := "not-so-secret-key"
	if dpv.ConfigInstance != nil {
		secretKey = dpv.ConfigInstance.Auth.DpvSecretKey
	}
	entropy := "f3ctRkFcqdgyXjSAleutn0UDx22/DZ8DlfmLNfHGtl8"
	innerKey := []byte(secretKey + entropy)
	outerKey := []byte(secretKey + entropy)
	for i := 0; i < len(innerKey); i++ {
		innerKey[i] = innerKey[i] ^ byte(0x36)
	}
	for i := 0; i < len(outerKey); i++ {
		outerKey[i] = outerKey[i] ^ byte(0x5c)
	}
	h := make([]byte, 32)
	d := sha3.NewShake128()
	d.Write(innerKey)
	d.Write([]byte(token))
	d.Read(h)
	d.Reset()
	d.Write(outerKey)
	d.Write(h)
	d.Read(h)
	return base64.RawURLEncoding.EncodeToString(h)
}

func IsStrongPassword(password string) bool {
	// minimum length: 8 characters
	// contains at least one character that is not a digit

	if len(password) < 8 {
		return false
	}
	for _, c := range password {
		if c < '0' || c > '9' {
			return true
		}
	}
	return false
}
