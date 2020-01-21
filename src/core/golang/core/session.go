package core

import (
	"time"
)

type UserSession struct {
	SessionId string `db:"session_id"`

	LastActiveTime time.Time `db:"last_active_time"`
	ExpirationTime time.Time `db:"expiration_time"`

	// TODO: Parse this out into browser/OS?
	UserAgent    string `db:"user_agent"`
	IP           string `db:"ip_address"`
	AccessToken  string `db:"access_token"`
	IdToken      string `db:"id_token"`
	RefreshToken string `db:"refresh_token"`

	UserId int64 `db:"user_id"`
}

// Other data that isn't necessarily something we want to call the "UserSesssion"
// but rather information that gets extracted out of the user session.
type UserSessionParsedData struct {
	CurrentUser    *User
	AccessibleOrgs []int32
	VerifiedEmail  bool
}