package security

import (
	"github.com/dvsekhvalnov/jose2go/base64url"
	"golang.org/x/crypto/sha3"
	"pkv/api/src/internal/dpv"
)

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
	return base64url.Encode(h)
}
