package core

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

type EncryptionAlgorithm int

const (
	// RSA PKCS1v15 with SHA256
	RSA256 EncryptionAlgorithm = iota
)

func VerifySignature(payload []byte, signature []byte, alg EncryptionAlgorithm, key interface{}) error {
	switch alg {
	case RSA256:
		return VerifyRSA256Signature(payload, signature, key)
	}
	return errors.New("Unknown encryption algorithm.")
}

func VerifyRSA256Signature(payload []byte, signature []byte, key interface{}) error {
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return errors.New("Bad RSA PublicKey")
	}

	hashed := sha256.Sum256(payload)
	return rsa.VerifyPKCS1v15(rsaKey, crypto.SHA256, hashed[:], signature)
}
