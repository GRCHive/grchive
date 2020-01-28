package core_test

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"math/rand"
	"reflect"
	"testing"
)

var NotRandomReader = rand.New(rand.NewSource(0))

func generateTestRsaPrivatePublicKey() *rsa.PrivateKey {
	key, _ := rsa.GenerateKey(NotRandomReader, 4096)
	return key
}

func TestVerifySignatureDispatch(t *testing.T) {
	var calledFunc = make(map[core.EncryptionAlgorithm]bool)
	var mockInterface = core.EncryptionInterface{
		VerifyRSA256Signature: func([]byte, []byte, interface{}) error {
			calledFunc[core.RSA256] = true
			return nil
		},
	}

	for _, test := range []core.EncryptionAlgorithm{
		core.RSA256,
	} {
		// Need to clear out the map every time to ensure that each test only
		// sets one key to true.
		calledFunc = make(map[core.EncryptionAlgorithm]bool)
		core.VerifySignature(mockInterface, []byte{}, []byte{}, test, nil)
		assert.True(t, calledFunc[test])
		assert.Equal(t, len(calledFunc), 1)
	}

	err := core.VerifySignature(mockInterface, []byte{}, []byte{}, core.MAX_ENC_ALG, nil)
	assert.Equal(t, err, core.ErrUnknownAlg)
}

func TestVerifyRSA256Signaure(t *testing.T) {
	// Not passing a rsa.PublicKey as the "key" should error out.
	err := core.VerifyRSA256Signature([]byte{}, []byte{}, "test")
	assert.Equal(t, err, core.ErrBadKey)

	// Now sign something manually and use the VerifyRSA256Signature function to verify.
	key := generateTestRsaPrivatePublicKey()

	message := []byte("TestVerifyRSA256Signaure!")
	hashed := sha256.Sum256(message)
	signature, _ := rsa.SignPKCS1v15(NotRandomReader, key, crypto.SHA256, hashed[:])
	err = core.VerifyRSA256Signature(message, signature, &key.PublicKey)
	assert.Nil(t, err)

	err = core.VerifyRSA256Signature([]byte("NotThMessage"), signature, &key.PublicKey)
	assert.NotNil(t, err)
}

func TestDefaultEncryptionInterface(t *testing.T) {
	// Technically using Pointer() here relies on undefined behavior
	// since Pointer() isn't guaranteed to return a uniquely identifying
	// value for a function but since we're not using closures here I think
	// it's OK.
	assert.Equal(
		t,
		reflect.ValueOf(core.DefaultEncryptionInterface.VerifyRSA256Signature).Pointer(),
		reflect.ValueOf(core.VerifyRSA256Signature).Pointer())
}
