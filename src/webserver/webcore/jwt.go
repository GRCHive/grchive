package webcore

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"strings"
	"time"
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

	// Marshaled data.
	Header struct {
		Algorithm string `json:"alg"`
		Kid       string `json:"kid"`
	}

	Payload struct {
		Iss    string   `json:"iss"`
		Cid    string   `json:"cid"`
		Exp    int64    `json:"exp"`
		Aud    string   `json:"aud"`
		Email  string   `json:"email"`
		Sub    string   `json:"sub"`
		Groups []string `json:"groups"`
	}
}

var base64EncodingNoPadding = base64.URLEncoding.WithPadding(base64.NoPadding)

func (this *RawJWT) PayloadToSign() []byte {
	strData := base64EncodingNoPadding.EncodeToString(this.RawHeader) + "." +
		base64EncodingNoPadding.EncodeToString(this.RawPayload)
	return []byte(strData)
}

func (this *RawJWT) verifyIss() error {
	if core.LoadEnvConfig().Login.BaseUrl != this.Payload.Iss {
		return errors.New("Iss mismatch.")
	}
	return nil
}

func (this *RawJWT) verifyCid() error {
	if core.LoadEnvConfig().Login.ClientId != this.Payload.Cid {
		return errors.New("Cid mismatch.")
	}
	return nil
}

func (this *RawJWT) verifyExp() error {
	expTime := time.Unix(this.Payload.Exp, 0)

	if time.Now().Sub(expTime).Seconds() > core.LoadEnvConfig().Login.TimeDriftLeewaySeconds {
		return errors.New("Invalid exp.")
	}
	return nil
}

func (this *RawJWT) verifyAud(isAccessToken bool) error {
	if isAccessToken {
		if core.LoadEnvConfig().Login.AuthAudience != this.Payload.Aud {
			return errors.New("Aud mismatch.")
		}

	} else {
		if core.LoadEnvConfig().Login.ClientId != this.Payload.Aud {
			return errors.New("Aud mismatch.")
		}
	}
	return nil
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

	err = json.Unmarshal(retData.RawPayload, &retData.Payload)
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

// Returns the decrypted payload along with an error (if there was one).
// The error is not nil if the token can't be verified or some other error appeared.
func (this JWTManager) VerifyJWT(input string, isAccessToken bool) (*RawJWT, error) {
	jwt, err := ReadRawJWTFromString(input)
	if err != nil {
		return nil, err
	}

	keys, err := this.GetKeysFromKid(jwt.Header.Kid)
	if err != nil {
		return nil, err
	}

	payload := jwt.PayloadToSign()
	for i := 0; i < len(keys); i++ {
		// Assume RSA256 for now. May want to expand to different options in the future.
		// That or use a JWT library. Heh.
		err = core.VerifySignature(payload, jwt.RawSignature, core.RSA256, keys[0])
		if err != nil {
			continue
		}

		err = jwt.verifyIss()
		if err != nil {
			return nil, err
		}
		err = jwt.verifyAud(isAccessToken)
		if err != nil {
			return nil, err
		}

		if isAccessToken {
			err = jwt.verifyCid()
			if err != nil {
				return nil, err
			}
		}

		err = jwt.verifyExp()
		if err != nil {
			return nil, err
		}

		return jwt, nil
	}
	return nil, err
}