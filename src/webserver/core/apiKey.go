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

func (key ApiKey) SecondsToExpiration() int {
	now := time.Now().UTC()
	return int(key.ExpirationDate.UTC().Sub(now).Seconds())
}

func (key ApiKey) IsExpired() bool {
	return key.SecondsToExpiration() <= 0
}
