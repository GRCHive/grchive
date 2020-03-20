package storage

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoHMAC       error = errors.New("No HMAC found.")
	ErrHMACMismatch       = errors.New("HMAC could not be verified.")
)

func computeHMACSHA512(data []byte, key []byte) ([]byte, error) {
	mac := hmac.New(sha512.New512_256, key)
	_, err := mac.Write(data)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func appendHMACSHA512(data []byte, key []byte) ([]byte, error) {
	mac, err := computeHMACSHA512(data, key)
	if err != nil {
		return nil, err
	}

	b64Data := base64.StdEncoding.EncodeToString(data)
	b64Mac := base64.StdEncoding.EncodeToString(mac)
	combined := fmt.Sprintf("%s.%s", b64Data, b64Mac)
	return []byte(combined), nil
}

func verifyHMACSHA512(data []byte, key []byte) ([]byte, error) {
	split := strings.Split(string(data), ".")
	if len(split) != 2 {
		return data, ErrNoHMAC
	}

	rawData, err := base64.StdEncoding.DecodeString(split[0])
	if err != nil {
		return nil, err
	}

	refMac, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return nil, err
	}

	testMac, err := computeHMACSHA512(rawData, key)
	if err != nil {
		return nil, err
	}

	if !hmac.Equal(testMac, refMac) {
		return rawData, ErrHMACMismatch
	}

	return rawData, nil
}
