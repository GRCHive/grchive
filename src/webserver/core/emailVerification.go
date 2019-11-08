package core

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type EmailVerification struct {
	UserId               int64    `db:"user_id"`
	Code                 string   `db:"code"`
	VerificationSent     NullTime `db:"verification_sent"`
	VerificationReceived NullTime `db:"verification_received"`
}

func CreateNewEmailVerification(user *User) EmailVerification {
	veri := EmailVerification{
		UserId:           user.Id,
		Code:             uuid.New().String(),
		VerificationSent: NullTime{sql.NullTime{time.Now(), true}},
	}
	return veri
}
