package core

import (
	"crypto/sha512"
	"encoding/hex"
	"time"
)

type RawApiKey string

type ApiKey struct {
	Id             int64     `db:"id"`
	HashedKey      string    `db:"hashed_api_key"`
	ExpirationDate time.Time `db:"expiration_date"`
	UserId         int64     `db:"user_id"`
}

// I don't think using SHA512 here is necessarily insecure
// since it's just being used to store Api keys. Realistically,
// if any attackers gets this far they probably already have
// access to whatever this API key could give them...
func (key RawApiKey) Hash() string {
	hash := sha512.Sum512([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (key ApiKey) Matches(rawKey RawApiKey) bool {
	hashedRawKey := rawKey.Hash()
	return key.HashedKey == hashedRawKey
}

func (key ApiKey) SecondsToExpiration(c Clock) int {
	now := c.Now().UTC()
	return int(key.ExpirationDate.UTC().Sub(now).Seconds())
}

func (key ApiKey) IsExpired(c Clock) bool {
	return key.SecondsToExpiration(c) <= 0
}

// Refresh if only 10 Minutes left.
func (key ApiKey) NeedsRefresh(c Clock) bool {
	return key.SecondsToExpiration(c) <= 10*60+1
}
