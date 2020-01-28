package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
)

func setupHasher() {
	core.EnvConfig = &core.EnvConfigData{
		HashId: &core.HashIdConfigData{
			MinLength: 8,
			Salt:      "7409d5d933caf02affce",
		},
	}

	// As long as we don't panic we're OK.
	core.InitializeHasher()
}

func TestInitializeHasher(t *testing.T) {
	setupHasher()
	tempHasher := core.Hasher

	// Calling InitializeHasher twice should not do anything.
	core.InitializeHasher()
	assert.Equal(t, tempHasher, core.Hasher)
}

func TestHashId(t *testing.T) {
	setupHasher()

	// Test that we can hash and reverse hash accurately (which is all that we care about).
	for _, ref := range []int64{
		10,
		0,
		512,
		9000,
	} {
		hash, err := core.HashId(ref)
		assert.Nil(t, err)

		test, err := core.ReverseHashId(hash)
		assert.Nil(t, err)

		assert.Equal(t, ref, test)
	}
}
