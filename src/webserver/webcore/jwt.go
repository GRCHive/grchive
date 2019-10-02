package webcore

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"strings"
)

// Deals with obtaining a JSON Web Key (JWK) (and re-obtaining if necessary) and using the
// cached keys to verify JSON Web Tokens (JWT).
type JWTImpl interface {
	// Maps key id to the public key. Apparently the kid -> key mapping need not be unique.
	// TODO: Support more than RSA?
	RetrieveKeys() (map[string][]*rsa.PublicKey, error)
}

type JWTManager struct {
	impl       JWTImpl
	cachedKeys map[string][]*rsa.PublicKey
}

type RawJWT struct {
	// Raw decoded string data.
	RawHeader    []byte
	RawPayload   []byte
	RawSignature []byte

	// Decoded header
	Header struct {
		Algorithm string `json:"alg"`
		Kid       string `json:"kid"`
	}
}

var base64EncodingNoPadding = base64.URLEncoding.WithPadding(base64.NoPadding)

func (this *RawJWT) PayloadToSign() []byte {
	strData := base64EncodingNoPadding.EncodeToString(this.RawHeader) + "." +
		base64EncodingNoPadding.EncodeToString(this.RawPayload)
	return []byte(strData)
}

func ReadRawJWTFromString(input string) (*RawJWT, error) {
	var err error
	data := strings.Split(input, ".")
	retData := &RawJWT{}

	retData.RawHeader, err = base64EncodingNoPadding.DecodeString(data[0])
	if err != nil {
		return nil, err
	}

	retData.RawPayload, err = base64EncodingNoPadding.DecodeString(data[1])
	if err != nil {
		return nil, err
	}

	retData.RawSignature, err = base64EncodingNoPadding.DecodeString(data[2])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(retData.RawHeader, &retData.Header)
	if err != nil {
		return nil, err
	}
	return retData, nil
}

func (this *JWTManager) GetKeysFromKid(kid string) ([]*rsa.PublicKey, error) {
	if key, ok := this.cachedKeys[kid]; ok {
		return key, nil
	}

	var err error
	this.cachedKeys, err = this.impl.RetrieveKeys()
	if err != nil {
		return nil, err
	}

	if key, ok := this.cachedKeys[kid]; ok {
		return key, nil
	}
	return nil, errors.New("Failed to find kid: " + kid)
}

// Returns not nil if the token can't be verified or some other error appeared.
func (this JWTManager) Verify(input string) error {
	jwt, err := ReadRawJWTFromString(input)
	if err != nil {
		return err
	}

	keys, err := this.GetKeysFromKid(jwt.Header.Kid)
	if err != nil {
		return err
	}

	payload := jwt.PayloadToSign()
	for i := 0; i < len(keys); i++ {
		// Assume RSA256.
		err = core.VerifySignature(payload, jwt.RawSignature, core.RSA256, keys[0])
		if err == nil {
			return nil
		}
	}
	return err
}
