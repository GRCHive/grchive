package core

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHexString(nBytes int) (string, error) {
	bytes := make([]byte, nBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", nil
	}
	return hex.EncodeToString(bytes), nil
}
