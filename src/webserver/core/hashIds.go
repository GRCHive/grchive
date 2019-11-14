package core

// This file deals with hashing using the method found in
// https://hashids.org/
import "github.com/speps/go-hashids"

var hasher *hashids.HashID = nil

func initializeHasher() {
	if hasher != nil {
		return
	}

	data := hashids.NewData()
	data.MinLength = EnvConfig.HashId.MinLength
	data.Salt = EnvConfig.HashId.Salt

	var err error
	hasher, err = hashids.NewWithData(data)
	if err != nil {
		panic("Failed to initialize hasher: " + err.Error())
	}
}

func HashId(c int64) (string, error) {
	return hasher.EncodeInt64([]int64{c})
}

func ReverseHashId(c string) (int64, error) {
	d, err := hasher.DecodeInt64WithError(c)
	if err != nil {
		return -1, err
	}

	return d[0], nil
}
