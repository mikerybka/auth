package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func newSessionToken() string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
