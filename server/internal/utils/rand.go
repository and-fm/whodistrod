package utils

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

func GenRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	encodedString := base32.StdEncoding.EncodeToString(b)
	encodedString = strings.TrimRight(encodedString, "=")

	return strings.ToLower(encodedString)
}
