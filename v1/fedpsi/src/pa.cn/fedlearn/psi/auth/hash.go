package auth

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
)

var (
	saltBytes = bytes.NewBufferString("PingAn-is-Great!").Bytes()
)

func HashPassword(password string) string {
	passwordBytes := []byte(password)

	sha512Hasher := sha512.New()

	passwordBytes = append(passwordBytes, saltBytes...)

	sha512Hasher.Write(passwordBytes)

	hashedPasswordBytes := sha512Hasher.Sum(nil)
	base64EncodedPasswordHash := base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}
