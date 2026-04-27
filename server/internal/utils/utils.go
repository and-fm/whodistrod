package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type J = map[string]interface{}

func Ctb() context.Context {
	return context.Background()
}

func HashSessionIdFromToken(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	sessionId := hex.EncodeToString(hash.Sum(nil))

	return sessionId
}
