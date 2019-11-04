package core

import (
	"crypto/sha512"
	"time"
)

type RawApiKey string

type ApiKey struct {
	Id             int64     `db:"id"`
	HashedKey      []byte    `db:"hashed_api_key"`
	Salt           string    `db:"salt"`
	ExpirationDate time.Time `db:"expiration_date"`
	UserId         int64     `db:"user_id"`
}

// I don't think using SHA512 here is necessarily insecure
// since it's just being used to store Api keys. Realistically,
// if any attackers gets this far they probably already have
// access to whatever this API key could give them...
func (key RawApiKey) HashWithSalt(salt string) []byte {
	hash := sha512.Sum512([]byte(string(key) + salt))
	return hash[:]
}

func (key ApiKey) Matches(rawKey RawApiKey) bool {
	hashedRawKey := rawKey.HashWithSalt(key.Salt)

	if key.HashedKey == nil || hashedRawKey == nil {
		return false
	}

	if len(key.HashedKey) != len(hashedRawKey) {
		return false
	}

	for i, b := range key.HashedKey {
		if b != hashedRawKey[i] {
			return false
		}
	}

	return true
}

func (key ApiKey) SecondsToExpiration() int {
	now := time.Now().UTC()
	return int(key.ExpirationDate.UTC().Sub(now).Seconds())
}

func (key ApiKey) IsExpired() bool {
	return key.SecondsToExpiration() <= 0
}
