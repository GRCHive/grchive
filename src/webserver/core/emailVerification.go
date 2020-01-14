package core

import (
	"database/sql"
)

type EmailVerification struct {
	UserId               int64    `db:"user_id"`
	Code                 string   `db:"code"`
	VerificationSent     NullTime `db:"verification_sent"`
	VerificationReceived NullTime `db:"verification_received"`
}

func (a EmailVerification) Equal(b EmailVerification) bool {
	return a.UserId == b.UserId && a.Code == b.Code &&
		a.VerificationSent.Equal(b.VerificationSent) &&
		a.VerificationReceived.Equal(b.VerificationReceived)
}

func CreateNewEmailVerification(user *User, u UuidGen, c Clock) EmailVerification {
	veri := EmailVerification{
		UserId:           user.Id,
		Code:             u.GenStr(),
		VerificationSent: NullTime{sql.NullTime{c.Now(), true}},
	}
	return veri
}
