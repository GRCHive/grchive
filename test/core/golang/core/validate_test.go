package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
)

func TestValidateEmailFormat(t *testing.T) {
	// Test cases source:
	// https://blogs.msdn.microsoft.com/testing123/2009/02/06/email-address-test-cases/
	for _, ref := range []struct {
		email string
		valid bool
	}{
		{`email@domain.com`, true},
		{`firstname.lastname@domain.com`, true},
		{`email@subdomain.domain.com`, true},
		{`firstname+lastname@domain.com`, true},
		{`email@123.123.123.123`, true},
		{`email@[123.123.123.123]`, true},
		{`"email"@domain.com`, true},
		{`1234567890@domain.com`, true},
		{`email@domain-one.com`, true},
		{`_______@domain.com`, true},
		{`email@domain.name`, true},
		{`email@domain.co.jp`, true},
		{`firstname-lastname@domain.com`, true},
		{`plainaddress`, false},
		{`#@%^%#$@#$@#.com`, false},
		{`@domain.com`, false},
		{`Joe Smith <email@domain.com>`, false},
		{`email.domain.com`, false},
		{`email@domain@domain.com`, false},
		{`.email@domain.com`, false},
		{`email.@domain.com`, false},
		{`email..email@domain.com`, false},
		{`あいうえお@domain.com`, false},
		{`email@domain.com (Joe Smith)`, false},
		{`email@domain`, false},
		{`email@-domain.com`, false},
		{`email@domain..com`, false},
	} {
		valid := core.ValidateEmailFormat(ref.email)
		assert.Equal(t, ref.valid, valid, ref.email)
	}

}
