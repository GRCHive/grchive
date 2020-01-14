package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"testing"
	"time"
)

type StubClock struct{}

func (c StubClock) Now() time.Time {
	utcLoc, _ := time.LoadLocation("UTC")
	return time.Date(2000, time.January, 1, 1, 1, 10, 1, utcLoc)
}

func TestRawApiKeyHash(t *testing.T) {
	for _, ref := range []struct {
		key  core.RawApiKey
		hash string
	}{
		{
			"ABC",
			"397118fdac8d83ad98813c50759c85b8c47565d8268bf10da483153b747a74743a58a90e85aa9f705ce6984ffc128db567489817e4092d050d8a1cc596ddc119",
		},
		{
			"",
			"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
	} {
		assert.Equal(t, ref.hash, ref.key.Hash())
	}
}

func TestApiKeyMatches(t *testing.T) {
	for _, ref := range []struct {
		key   core.ApiKey
		raw   core.RawApiKey
		match bool
	}{
		{
			key: core.ApiKey{
				HashedKey: "397118fdac8d83ad98813c50759c85b8c47565d8268bf10da483153b747a74743a58a90e85aa9f705ce6984ffc128db567489817e4092d050d8a1cc596ddc119",
			},
			raw:   "ABC",
			match: true,
		},
		{
			key: core.ApiKey{
				HashedKey: "abc",
			},
			raw:   "ABC",
			match: false,
		},
	} {
		assert.Equal(t, ref.match, ref.key.Matches(ref.raw))
	}
}

func TestSecondsToExpiration(t *testing.T) {
	utcLoc, _ := time.LoadLocation("UTC")

	for _, ref := range []struct {
		key       core.ApiKey
		clock     StubClock
		seconds   int
		isExpired bool
	}{
		{
			key: core.ApiKey{
				ExpirationDate: time.Date(2000, time.January, 1, 1, 1, 20, 1, utcLoc),
			},
			clock:     StubClock{},
			seconds:   10,
			isExpired: false,
		},
		{
			key: core.ApiKey{
				ExpirationDate: time.Date(2000, time.January, 1, 1, 1, 10, 1, utcLoc),
			},
			clock:     StubClock{},
			seconds:   0,
			isExpired: true,
		},
		{
			key: core.ApiKey{
				ExpirationDate: time.Date(2000, time.January, 1, 1, 1, 0, 1, utcLoc),
			},
			clock:     StubClock{},
			seconds:   -10,
			isExpired: true,
		},
	} {
		assert.Equal(t, ref.seconds, ref.key.SecondsToExpiration(ref.clock))
		assert.Equal(t, ref.isExpired, ref.key.IsExpired(ref.clock))
	}
}
