package webcore

import (
	"encoding/hex"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/vault_api"
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
