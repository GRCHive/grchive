package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
	"time"
)

type StubClock struct{}

func (c StubClock) Now() time.Time {
	utcLoc, _ := time.LoadLocation("UTC")
	return time.Date(2000, time.January, 1, 1, 1, 10, 1, utcLoc)
}

type StubUuidGen struct{}

func (g StubUuidGen) GenStr() string {
	return "ABCDEF"
}

func TestEmailVerificationEqual(t *testing.T) {
	utcLoc, _ := time.LoadLocation("UTC")
	t1 := time.Date(2000, time.January, 1, 1, 1, 10, 1, utcLoc)
	t2 := time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc)

	for _, ref := range []struct {
		a     core.EmailVerification
		b     core.EmailVerification
		match bool
	}{
		{
			a: core.EmailVerification{
				UserId: 525,
				Code:   "ABC",
			},
			b: core.EmailVerification{
				UserId: 525,
				Code:   "ABC",
			},
			match: true,
		},
		{
			a: core.EmailVerification{
				UserId:           525,
				Code:             "ABC",
				VerificationSent: core.CreateNullTime(t1),
			},
			b: core.EmailVerification{
				UserId:           525,
				Code:             "ABC",
				VerificationSent: core.CreateNullTime(t1),
			},
			match: true,
		},
		{
			a: core.EmailVerification{
				UserId:               525,
				Code:                 "ABC",
				VerificationReceived: core.CreateNullTime(t2),
			},
			b: core.EmailVerification{
				UserId:               525,
				Code:                 "ABC",
				VerificationReceived: core.CreateNullTime(t2),
			},
			match: true,
		},
		{
			a: core.EmailVerification{
				UserId: 526,
				Code:   "ABC",
			},
			b: core.EmailVerification{
				UserId: 525,
				Code:   "ABC",
			},
			match: false,
		},
		{
			a: core.EmailVerification{
				UserId: 525,
				Code:   "ABCD",
			},
			b: core.EmailVerification{
				UserId: 525,
				Code:   "ABC",
			},
			match: false,
		},
		{
			a: core.EmailVerification{
				UserId:           525,
				Code:             "ABC",
				VerificationSent: core.CreateNullTime(t1),
			},
			b: core.EmailVerification{
				UserId:           525,
				Code:             "ABC",
				VerificationSent: core.CreateNullTime(t2),
			},
			match: false,
		},
		{
			a: core.EmailVerification{
				UserId:           525,
				Code:             "ABC",
				VerificationSent: core.CreateNullTime(t1),
			},
			b: core.EmailVerification{
				UserId: 525,
				Code:   "ABC",
			},
			match: false,
		},
		{
			a: core.EmailVerification{
				UserId:               525,
				Code:                 "ABC",
				VerificationReceived: core.CreateNullTime(t1),
			},
			b: core.EmailVerification{
				UserId:               525,
				Code:                 "ABC",
				VerificationReceived: core.CreateNullTime(t2),
			},
			match: false,
		},
		{
			a: core.EmailVerification{
				UserId:               525,
				Code:                 "ABC",
				VerificationReceived: core.CreateNullTime(t1),
			},
			b: core.EmailVerification{
				UserId: 525,
				Code:   "ABC",
			},
			match: false,
		},
	} {
		assert.Equal(t, ref.match, ref.a.Equal(ref.b))
	}
}

func TestCreateNewEmailVerification(t *testing.T) {
	clock := StubClock{}
	uuid := StubUuidGen{}

	for _, ref := range []struct {
		user core.User
		veri core.EmailVerification
	}{
		{
			user: core.User{
				Id: 256,
			},
			veri: core.EmailVerification{
				UserId:           256,
				Code:             uuid.GenStr(),
				VerificationSent: core.CreateNullTime(clock.Now()),
			},
		},
		{
			user: core.User{
				Id: 5555,
			},
			veri: core.EmailVerification{
				UserId:           5555,
				Code:             uuid.GenStr(),
				VerificationSent: core.CreateNullTime(clock.Now()),
			},
		},
	} {
		assert.True(t, ref.veri.Equal(core.CreateNewEmailVerification(&ref.user, uuid, clock)))
	}
}
