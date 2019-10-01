package core

import (
	"time"
)

type UserSession struct {
	SessionId      string
	UserId         string
	LastActiveTime time.Time
	ExpirationTime time.Time
	Browser        string
	Location       string
}
