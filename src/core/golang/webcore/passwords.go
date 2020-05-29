package webcore

import (
	"encoding/hex"
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/vault_api"
	"strings"
)

const numBytesForSalt int = 64
const pwTransitPath string = "passwords"

// Note: encrypted, not hashed.
// Returns: encrypted, salt, error
func CreateSaltedEncryptedPassword(rawPassword string) (string, string, error) {
	salt, err := core.RandomHexString(numBytesForSalt)
	if err != nil {
		return "", "", err
	}

	basePw := hex.EncodeToString([]byte(rawPassword))
	fullPw := basePw + "." + salt
	encPwBytes, err := vault.TransitEncrypt(pwTransitPath, []byte(fullPw))
	if err != nil {
		return "", "", err
	}

	return string(encPwBytes), salt, nil
}

func DecryptSaltedEncryptedPassword(encPassword string, salt string) (string, error) {
	decPwBytes, err := vault.TransitDecrypt(pwTransitPath, []byte(encPassword))
	if err != nil {
		return "", err
	}

	splitData := strings.Split(string(decPwBytes), ".")
	if len(splitData) == 0 {
		return "", errors.New("Can't find encrypted password delimiter.")
	}

	// The last element should be the salt.
	if splitData[len(splitData)-1] != salt {
		return "", errors.New("Salt mismatch.")
	}

	decodePw, err := hex.DecodeString(strings.Join(splitData[:len(splitData)-1], "."))
	if err != nil {
		return "", err
	}

	return string(decodePw), nil
}

func CreateEncryptedPassword(rawPassword string) (string, error) {
	basePw := hex.EncodeToString([]byte(rawPassword))
	encPwBytes, err := vault.TransitEncrypt(pwTransitPath, []byte(basePw))
	if err != nil {
		return "", err
	}

	return string(encPwBytes), nil
}

func DecryptEncryptedPassword(encPassword string) (string, error) {
	decPwBytes, err := vault.TransitDecrypt(pwTransitPath, []byte(encPassword))
	if err != nil {
		return "", err
	}

	decodePw, err := hex.DecodeString(string(decPwBytes))
	if err != nil {
		return "", err
	}

	return string(decodePw), nil
}
