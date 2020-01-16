package core

import (
	"errors"
)

var ErrorUnauthorized = errors.New("Unauthorized access")

// Obtains the error string. If the error is nil, returns "No Error."
func ErrorString(err error) string {
	if err == nil {
		return "No Error"
	}
	return err.Error()
}
