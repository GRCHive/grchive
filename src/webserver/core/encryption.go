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
	MAX_ENC_ALG
)

var (
	ErrUnknownAlg error = errors.New("Unknown encryption algorithm")
	ErrBadKey           = errors.New("Bad key")
)

type VerifySignatureFunc func([]byte, []byte, interface{}) error

type EncryptionInterface struct {
	VerifyRSA256Signature VerifySignatureFunc
}

func VerifySignature(enc EncryptionInterface, payload []byte, signature []byte, alg EncryptionAlgorithm, key interface{}) error {
	switch alg {
	case RSA256:
		return enc.VerifyRSA256Signature(payload, signature, key)
	}
	return ErrUnknownAlg
}

func VerifyRSA256Signature(payload []byte, signature []byte, key interface{}) error {
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrBadKey
	}

	hashed := sha256.Sum256(payload)
	return rsa.VerifyPKCS1v15(rsaKey, crypto.SHA256, hashed[:], signature)
}

var DefaultEncryptionInterface = EncryptionInterface{
	VerifyRSA256Signature: VerifyRSA256Signature,
}
